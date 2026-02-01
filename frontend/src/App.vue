<template>
  <div id="app">
    <div class="app-container">
      <Sidebar
        :groups="groups"
        :selected-group="selectedGroup"
        :system-host-path="systemHostPath"
        :search-query="searchQuery"
        @select-group="selectGroup"
        @toggle-status="toggleGroupStatus"
        @delete-group="deleteGroup"
        @open-add-modal="showAddGroupModal = true"
        @update:search-query="searchQuery = $event"
        @select-system-host="selectSystemHost"
        @restore-system-hosts="restoreSystemHosts"
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
  RestoreData,
  RestoreRawSystemHosts,
  HasRawHostsBackup
} from '../wailsjs/go/main/App'

import { ElMessageBox, ElNotification } from 'element-plus';

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
      } catch (error) {
        this.showMessage(`Failed to load host groups: ${error}`, 'error')
      }
    },

    selectGroup(group) {
      if (this.isDirty) {
        // 使用 Element Plus 的 MessageBox 替换 confirm
        ElMessageBox.confirm(
          'You have unsaved changes. Are you sure you want to switch groups?',
          'Confirmation',
          {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            type: 'warning',
            draggable: true
          }
        ).then(() => {
          // 用户点击了确认
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
        }).catch(() => {
          // 用户点击了取消，什么都不做
        });
      } else {
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
      }
    },

    async toggleGroupStatus(group) {
      const newStatus = !group.enabled
      try {
        await ToggleHostGroup(group.id, newStatus)
        
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
      // 使用 Element Plus 的 MessageBox 替换 confirm
      try {
        await ElMessageBox.confirm(
          'Are you sure you want to delete this group?',
          'Confirm Delete',
          {
            confirmButtonText: 'Delete',
            cancelButtonText: 'Cancel',
            type: 'warning',
            draggable: true
          }
        );
        
        // 用户点击了确认
        await DeleteHostGroup(groupId)
        await this.loadHostGroups()
        if (this.selectedGroup && this.selectedGroup.id === groupId) {
          this.selectedGroup = null
          this.editingGroup = {}
          this.remoteContentPreview = ''
        }
        this.showMessage('Group deleted successfully', 'success')
      } catch (error) {
        // 用户点击了取消或关闭对话框
        if (error !== 'cancel') {
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
        this.showMessage(`App data backup completed successfully: ${result}`, 'success');
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

        // 构建备份列表的HTML字符串
        let backupListHtml = `<div style="max-height: 400px; overflow-y: auto;">
          <p style="margin-bottom: 15px;">请选择要恢复的备份文件：</p>`;
        
        backups.forEach((backup, index) => {
          backupListHtml += `
            <div style="display: flex; justify-content: space-between; align-items: center; margin: 8px 0; padding: 12px; border: 1px solid #dcdfe6; border-radius: 4px; background-color: #fafafa;">
              <span style="flex: 1; word-break: break-all; margin-right: 10px; font-size: 14px;">${backup}</span>
              <button id="restore-btn-${index}" style="padding: 6px 12px; background-color: #409eff; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 13px;">恢复</button>
            </div>`;
        });
        
        backupListHtml += '</div>';
        
        // 创建临时div来包含HTML内容
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = backupListHtml;
        
        // 显示包含备份列表的对话框
        const msgBox = ElMessageBox({
          title: 'Restore Backup',
          dangerouslyUseHTMLString: true,
          message: tempDiv.innerHTML,
          showCancelButton: false,
          showConfirmButton: true,
          confirmButtonText: 'Close',
          closeOnClickModal: false,
          closeOnPressEscape: true,
          customStyle: { 
            width: '600px',
            maxHeight: '500px'
          }
        });
        
        // 等待对话框显示后再绑定事件
        setTimeout(() => {
          backups.forEach((backup, index) => {
            const button = document.getElementById(`restore-btn-${index}`);
            if (button) {
              const handleClick = async () => {
                try {
                  // 关闭当前对话框
                  ElMessageBox.close();
                  
                  // 确认恢复操作
                  await ElMessageBox.confirm(
                    `Are you sure you want to restore from ${backup}?\n\nThis will overwrite your current data!`,
                    'Confirm Restore',
                    {
                      confirmButtonText: 'Restore',
                      cancelButtonText: 'Cancel',
                      type: 'warning',
                      draggable: true
                    }
                  );

                  console.log('Starting restore...');
                  // 执行恢复
                  await RestoreData(backup);
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
                  this.showMessage(`Restored ${this.groups.length} groups (${enabledCount} enabled) from ${backup}`, 'success');
                } catch (error) {
                  if (error !== 'cancel' && error?.type !== 'cancel') {
                    this.showMessage(`Failed to restore backup: ${error}`, 'error');
                  }
                  // 重新显示备份列表
                  setTimeout(() => {
                    this.restoreBackup();
                  }, 100);
                }
              };
              
              button.addEventListener('click', handleClick);
            }
          });
        }, 100);

      } catch (error) {
        if (error !== 'cancel' && error?.type !== 'cancel') {
          this.showMessage(`Failed to restore backup: ${error}`, 'error');
        }
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

    async restoreSystemHosts() {
      try {
        // 检查是否存在原始hosts备份文件
        const hasBackup = await HasRawHostsBackup();
        
        console.log('Has raw hosts backup:', hasBackup); // 调试信息
        
        if (!hasBackup) {
          this.showMessage('No system hosts backup found! Expected file: raw_hosts_backup.txt', 'error');
          return;
        }
        
        // 使用 Element Plus 的 MessageBox 替换 confirm
        await ElMessageBox.confirm(
          'Are you sure you want to restore the system hosts file to its original state using the backup? This will revert all changes made by Ghost Host Manager.',
          'Confirm Restore System Hosts',
          {
            confirmButtonText: 'Restore',
            cancelButtonText: 'Cancel',
            type: 'warning',
            draggable: true
          }
        );
        
        // 执行恢复操作
        await RestoreRawSystemHosts('raw_hosts_backup.txt');
        
        // 自动禁用所有其他host分组
        await this.disableAllHostGroups();
        
        this.showMessage('System hosts restored successfully! All host groups have been disabled.', 'success');
      } catch (error) {
        if (error !== 'cancel' && error?.type !== 'cancel') {
          console.error('Failed to restore system hosts:', error);
          this.showMessage(`Failed to restore system hosts: ${error}`, 'error');
        }
      }
    },
    
    // 自动禁用所有host分组
    async disableAllHostGroups() {
      try {
        // 获取所有host分组
        const allGroups = await GetHostGroups();
        
        // 遍历所有启用的分组并禁用它们
        for (const group of allGroups) {
          if (group.enabled) {
            await ToggleHostGroup(group.id, false);
          }
        }
        
        // 重新加载host分组以更新UI
        await this.loadHostGroups();
        
        // 应用更改到系统hosts文件
        await this.applyHosts();
      } catch (error) {
        console.error('Failed to disable host groups:', error);
        this.showMessage(`Failed to disable host groups: ${error}`, 'error');
      }
    },

    async refreshSystemHost() {
      try {
        this.systemHostContent = await GetSystemHostsContent()
        if (this.selectedGroup && this.selectedGroup.id === 'system-host') {
          this.selectedGroup.content = this.systemHostContent
        }
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
      // 使用 Element Plus 的 Notification 替换自定义消息
      switch (type) {
        case 'success':
          ElNotification({
            title: 'Success',
            message: text,
            type: 'success',
            duration: 3000,
            position: 'top-right'
          });
          break;
        case 'error':
          ElNotification({
            title: 'Error',
            message: text,
            type: 'error',
            duration: 5000,
            position: 'top-right'
          });
          break;
        case 'info':
          ElNotification({
            title: 'Info',
            message: text,
            type: 'info',
            duration: 3000,
            position: 'top-right'
          });
          break;
        default:
          ElNotification({
            title: 'Message',
            message: text,
            type: 'info',
            duration: 3000,
            position: 'top-right'
          });
      }
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
