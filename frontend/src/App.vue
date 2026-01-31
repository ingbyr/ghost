<template>
  <div id="app">
    <div class="app-container">
      <Sidebar 
        :groups="groups"
        :selected-group="selectedGroup"
        :system-host-path="systemHostPath"
        :search-query="searchQuery"
        @select-group="selectGroup"
        @select-system-host="selectSystemHost"
        @toggle-status="toggleGroupStatus"
        @delete-group="deleteGroup"
        @open-add-modal="showAddGroupModal = true"
        @update:search-query="searchQuery = $event"
      />
      <MainPanel
        :selected-group="selectedGroup"
        :editing-group="editingGroup"
        :system-host-path="systemHostPath"
        :system-host-content="systemHostContent"
        :remote-content-preview="remoteContentPreview"
        :is-fetching-remote="isFetchingRemote"
        @save-group="saveGroup"
        @cancel-edit="cancelEdit"
        @fetch-remote-content="fetchRemoteContent"
        @mark-as-dirty="markAsDirty"
        @refresh-system-host="refreshSystemHost"
        @apply-hosts="applyHosts"
      />
    </div>
    
    <ActionBar
      @refresh-remote="refreshRemote"
      @refresh-list="loadHostGroups"
      @backup-now="backupNow"
      @restore-backup="restoreBackup"
    />
    
    <AddGroupModal
      :show-modal="showAddGroupModal"
      :new-group="newGroup"
      @close-modal="closeAddGroupModal"
      @add-group="addGroup"
      @update:newGroup="updateNewGroup"
    />
    
    <!-- 消息提示 -->
    <div v-if="message" :class="['message', messageType]">
      {{ message }}
    </div>
  </div>
</template>

<script>
import { 
  GetHostGroups, 
  AddHostGroup, 
  UpdateHostGroup, 
  DeleteHostGroup, 
  ToggleHostGroup, 
  ApplyHosts, 
  RefreshRemoteGroups,
  GetSystemHostsContent,
  GetSystemHostPath,
  GetRemoteContent,
  RefreshRemoteGroup,
  BackupConfig,
  CreateSystemHostsBackup,
  BackupAppAndSystemHosts,
  ListDataBackups,
  RestoreData
} from '../wailsjs/go/main/App'

import Sidebar from './components/Sidebar.vue';
import MainPanel from './components/MainPanel.vue';
import ActionBar from './components/ActionBar.vue';
import AddGroupModal from './components/AddGroupModal.vue';

export default {
  name: 'App',
  components: {
    Sidebar,
    MainPanel,
    ActionBar,
    AddGroupModal
  },
  data() {
    return {
      groups: [],
      selectedGroup: null,
      editingGroup: {},
      newGroup: {
        name: '',
        description: '',
        content: '',
        enabled: false,
        isRemote: false,
        url: ''
      },
      showAddGroupModal: false,
      searchQuery: '',
      message: '',
      messageType: '', // 'success', 'error', 'info'
      isDirty: false,
      systemHostPath: '',
      systemHostContent: '',
      remoteContentPreview: '',
      isRefreshingRemote: false,
      isFetchingRemote: false
    }
  },
  async mounted() {
    await this.loadHostGroups()
    await this.selectSystemHost()
  },
  methods: {
    async loadHostGroups() {
      try {
        console.log('Loading host groups...');
        this.groups = await GetHostGroups()
        console.log('Host groups loaded:', this.groups.length, 'groups');
        if (this.selectedGroup) {
          // 如果当前选中的是系统hosts文件，保持选择不变
          if (this.selectedGroup.id === 'system-host') {
            // 但仍然更新系统hosts内容
            await this.selectSystemHost();
          } else {
            // 更新选中的组（以防数据有变化）
            const updatedGroup = this.groups.find(g => g.id === this.selectedGroup.id);
            if (updatedGroup) {
              this.selectedGroup = updatedGroup;
              this.editingGroup = { ...updatedGroup };
              this.isDirty = false;
            } else {
              // 如果找不到对应的组，清空选择
              this.selectedGroup = null;
              this.editingGroup = {};
              this.isDirty = false;
            }
          }
        }
        this.showMessage('Groups loaded successfully', 'info')
      } catch (error) {
        this.showMessage(`Failed to load host groups: ${error}`, 'error')
      }
    },

    selectGroup(group) {
      if (this.isDirty) {
        if (!confirm('You have unsaved changes. Are you sure you want to switch groups?')) {
          return
        }
      }
      this.selectedGroup = group
      this.editingGroup = { ...group }
      this.isDirty = false
      
      // 如果是远程组，将组内容设置为预览内容
      if (group.isRemote) {
        this.remoteContentPreview = group.content || '';
      } else {
        // 如果是非远程组，清空预览内容
        this.remoteContentPreview = '';
      }
    },

    async toggleGroupStatus(group) {
      const newStatus = !group.enabled
      try {
        await ToggleHostGroup(group.id, newStatus)
        this.showMessage(`Group ${newStatus ? 'enabled' : 'disabled'}`, 'success')
        
        // 重新加载组数据以确保UI状态一致
        await this.loadHostGroups()
        
        // 启用或禁用组后，应用所有启用的Hosts
        await this.applyHosts()
      } catch (error) {
        this.showMessage(`Failed to toggle group: ${error}`, 'error')
      }
    },
    
    async refreshSpecificRemoteGroup(groupId) {
      try {
        // 通过后端API刷新指定的远程组内容
        await RefreshRemoteGroup(groupId);
        
        // 重新加载所有组以获取更新后的内容
        await this.loadHostGroups();
        
        // 如果当前选中的是这个组，更新本地状态
        if (this.selectedGroup && this.selectedGroup.id === groupId) {
          this.selectedGroup = this.groups.find(g => g.id === groupId);
          this.editingGroup = { ...this.selectedGroup };
          this.remoteContentPreview = this.selectedGroup.content;
        }
        
        // 如果组是启用状态，则应用到系统
        if (this.selectedGroup && this.selectedGroup.enabled) {
          await this.applyHosts();
        }
        
        this.showMessage('Remote group content updated successfully', 'success');
      } catch (error) {
        this.showMessage(`Failed to refresh remote group: ${error}`, 'error')
      }
    },

    async addGroup(newGroupData) {
      if (!newGroupData.name) {
        this.showMessage('Name is required', 'error')
        return
      }

      try {
        await AddHostGroup({ ...newGroupData })
        this.newGroup = {
          name: '',
          description: '',
          content: '',
          enabled: false,
          isRemote: false,
          url: ''
        }
        this.closeAddGroupModal()
        await this.loadHostGroups()
        this.showMessage('Group added successfully', 'success')
      } catch (error) {
        this.showMessage(`Failed to add group: ${error}`, 'error')
      }
    },

    async deleteGroup(groupId) {
      if (confirm('Are you sure you want to delete this group?')) {
        try {
          await DeleteHostGroup(groupId)
          await this.loadHostGroups()
          if (this.selectedGroup && this.selectedGroup.id === groupId) {
            this.selectedGroup = null
            this.editingGroup = {}
            this.isDirty = false
          }
          this.showMessage('Group deleted successfully', 'success')
        } catch (error) {
          this.showMessage(`Failed to delete group: ${error}`, 'error')
        }
      }
    },

    async saveGroup(editingGroupData) {
      try {
        await UpdateHostGroup(editingGroupData)
        // 更新主列表中的组
        const index = this.groups.findIndex(g => g.id === editingGroupData.id)
        if (index !== -1) {
          this.groups[index] = { ...editingGroupData }
          this.selectedGroup = { ...editingGroupData }
        }
        this.isDirty = false
        this.showMessage('Group updated successfully', 'success')
        
        // 如果该组已启用，则自动应用Hosts
        if (editingGroupData.enabled) {
          await this.applyHosts();
        }
      } catch (error) {
        this.showMessage(`Failed to update group: ${error}`, 'error')
      }
    },

    cancelEdit() {
      if (this.selectedGroup) {
        this.editingGroup = { ...this.selectedGroup }
        this.isDirty = false
      }
    },

    markAsDirty() {
      this.isDirty = true
    },

    async applyHosts() {
      try {
        await ApplyHosts()
        
        // 刷新系统hosts内容，如果当前显示的是系统hosts
        if (this.selectedGroup && this.selectedGroup.id === 'system-host') {
          await this.refreshSystemHost()
        }
      } catch (error) {
        this.showMessage(`Failed to apply hosts: ${error}`, 'error')
      }
    },

    async refreshRemote() {
      try {
        // 记录刷新前启用的远程组
        const previouslyEnabledGroups = this.groups.filter(g => g.enabled && g.isRemote);
        
        await RefreshRemoteGroups()
        await this.loadHostGroups()
        
        // 检查是否有任何启用的组，如果有，则应用到系统
        const hasEnabledGroups = this.groups.some(g => g.enabled);
        if (hasEnabledGroups) {
          await this.applyHosts();
        }
        
        this.showMessage('Remote groups refreshed', 'success')
      } catch (error) {
        this.showMessage(`Failed to refresh remote groups: ${error}`, 'error')
      }
    },
    
    async backupNow() {
      try {
        const result = await BackupAppAndSystemHosts();
        this.showMessage(`Backup completed successfully: ${result}`, 'success');
      } catch (error) {
        this.showMessage(`Failed to perform backup: ${error}`, 'error');
      }
    },

    async restoreBackup() {
      try {
        // 获取备份列表
        const backups = await ListDataBackups();
        if (backups.length === 0) {
          this.showMessage('No backups available', 'error');
          return;
        }

        // 格式化备份列表供用户选择
        const backupOptions = backups.map((backup, index) => {
          return `${index + 1}. ${backup}`;
        }).join('\n');

        // 显示选择对话框
        const selectedIndex = prompt(
          `Select a backup to restore (1-${backups.length}):\n\n${backupOptions}\n\nWarning: This will overwrite current data!`
        );

        if (selectedIndex === null) {
          return; // 用户取消
        }

        const index = parseInt(selectedIndex) - 1;
        if (isNaN(index) || index < 0 || index >= backups.length) {
          this.showMessage('Invalid selection', 'error');
          return;
        }

        const selectedBackup = backups[index];
        console.log('Selected backup:', selectedBackup);

        // 确认恢复
        if (!confirm(`Are you sure you want to restore from ${selectedBackup}?\n\nThis will overwrite your current data!`)) {
          console.log('User cancelled restore');
          return;
        }

        console.log('Starting restore...');
        // 执行恢复
        await RestoreData(selectedBackup);
        console.log('RestoreData completed');
        
        // 重置选中状态，避免状态混乱
        this.selectedGroup = null;
        this.editingGroup = {};
        this.isDirty = false;
        
        // 清空现有分组，强制 Vue 重新渲染
        this.groups = [];
        await this.$nextTick();
        
        // 刷新整个列表以显示恢复后的分组和启用状态
        await this.loadHostGroups();
        
        // 等待 Vue 更新 DOM
        await this.$nextTick();
        
        // 调试：输出加载的分组启用状态
        console.log('Loaded groups after restore:', this.groups.map(g => ({ id: g.id, name: g.name, enabled: g.enabled })));
        
        // 清空 Ghost 管理的系统 host 部分，应用当前启用的 host 分组
        await this.applyHosts();
        
        // 刷新系统 hosts 显示（因为 applyHosts 修改了系统 hosts）
        await this.refreshSystemHost();

        const enabledCount = this.groups.filter(g => g.enabled).length;
        this.showMessage(`Restored ${this.groups.length} groups (${enabledCount} enabled) from ${selectedBackup}`, 'success');
      } catch (error) {
        this.showMessage(`Failed to restore backup: ${error}`, 'error');
      }
    },

    closeAddGroupModal() {
      this.showAddGroupModal = false
      this.newGroup = {
        name: '',
        description: '',
        content: '',
        enabled: false,
        isRemote: false,
        url: '',
        typeSelected: false
      }
    },

    updateNewGroup(data) {
      this.newGroup = { ...data };
    },

    async selectSystemHost() {
      try {
        // 获取系统host路径和内容（每次都重新获取最新内容）
        this.systemHostPath = await GetSystemHostPath()
        this.systemHostContent = await GetSystemHostsContent()
        
        // 设置选中系统host
        this.selectedGroup = {
          id: 'system-host',
          name: 'System Host File',
          description: this.systemHostPath,
          content: this.systemHostContent,
          enabled: false,
          isRemote: false,
          createdAt: null,
          updatedAt: null
        }
        
        // 清除编辑状态
        this.editingGroup = {}
        this.isDirty = false
      } catch (error) {
        this.showMessage(`Failed to load system host file: ${error}`, 'error')
      }
    },

    async refreshSystemHost() {
      try {
        this.systemHostContent = await GetSystemHostsContent()
        if (this.selectedGroup && this.selectedGroup.id === 'system-host') {
          this.selectedGroup.content = this.systemHostContent
        }
        this.showMessage('System host file refreshed', 'info')
      } catch (error) {
        this.showMessage(`Failed to refresh system host file: ${error}`, 'error')
      }
    },

    async refreshRemoteContent() {
      if (!this.selectedGroup || !this.selectedGroup.isRemote || !this.selectedGroup.url) {
        this.showMessage('No remote URL configured for this group', 'error')
        return
      }

      // 保存原始的启用状态
      const wasEnabled = this.selectedGroup.enabled;
      
      this.isRefreshingRemote = true
      try {
        // 通过后端API刷新远程组内容，这会更新数据库中的内容
        await RefreshRemoteGroup(this.selectedGroup.id);
        
        // 重新加载所有组以获取更新后的内容
        await this.loadHostGroups();
        
        // 找到更新后的组并更新本地状态
        const updatedGroup = this.groups.find(g => g.id === this.selectedGroup.id);
        if (updatedGroup) {
          this.selectedGroup = updatedGroup;
          this.editingGroup = { ...updatedGroup };
          this.remoteContentPreview = updatedGroup.content;
        }
        
        // 如果组是启用状态，则应用到系统
        if (updatedGroup && updatedGroup.enabled) {
          await this.applyHosts();
        }
        
        this.showMessage('Remote content updated successfully', 'success')
      } catch (error) {
        this.showMessage(`Failed to update remote content: ${error}`, 'error')
      } finally {
        this.isRefreshingRemote = false
      }
    },

    async fetchRemoteContent() {
      if (!this.editingGroup.url) {
        this.showMessage('Please enter a URL first', 'error');
        return;
      }

      this.isFetchingRemote = true;
      try {
        // 从指定URL获取远程内容
        this.remoteContentPreview = await GetRemoteContent(this.editingGroup.url);
        
        // 如果当前编辑的是一个已保存的远程组，则更新其内容
        if (this.editingGroup.id && this.editingGroup.isRemote) {
          // 更新编辑组的content
          this.editingGroup.content = this.remoteContentPreview;
          
          // 更新选中组的content
          if (this.selectedGroup && this.selectedGroup.id === this.editingGroup.id) {
            this.selectedGroup.content = this.remoteContentPreview;
          }
          
          // 更新后保存到后端
          await UpdateHostGroup(this.editingGroup);
          
          // 重新加载所有组以确保数据一致性
          await this.loadHostGroups();
        }
        
        this.showMessage('Remote content fetched and saved successfully', 'success');
        
        // 如果当前选中的组已启用，则自动应用Hosts
        if (this.selectedGroup && this.selectedGroup.enabled) {
          await this.applyHosts();
        }
      } catch (error) {
        this.showMessage(`Failed to fetch remote content: ${error}`, 'error');
      } finally {
        this.isFetchingRemote = false;
      }
    },

    showMessage(text, type) {
      this.message = text
      this.messageType = type
      setTimeout(() => {
        this.message = ''
        this.messageType = ''
      }, 3000)
    },

    onSearchChanged(searchQuery) {
      // 处理搜索查询变化（如果需要额外逻辑）
    }
  }
}
</script>

<style>
#app {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.message {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 20px;
  border-radius: 4px;
  color: white;
  z-index: 2000;
  max-width: 400px;
  word-wrap: break-word;
}

.message.success {
  background-color: #28a745;
}

.message.error {
  background-color: #dc3545;
}

.message.info {
  background-color: #17a2b8;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
