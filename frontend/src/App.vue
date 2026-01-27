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
      />
    </div>
    
    <ActionBar
      @apply-hosts="applyHosts"
      @refresh-remote="refreshRemote"
      @refresh-list="loadHostGroups"
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
  RefreshRemoteGroup
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
        this.groups = await GetHostGroups()
        if (this.selectedGroup) {
          // 更新选中的组（以防数据有变化）
          this.selectedGroup = this.groups.find(g => g.id === this.selectedGroup.id) || null
          if (this.selectedGroup) {
            this.editingGroup = { ...this.selectedGroup }
            this.isDirty = false
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

    toggleGroupStatus(group) {
      const newStatus = !group.enabled
      ToggleHostGroup(group.id, newStatus)
        .then(() => {
          group.enabled = newStatus
          this.showMessage(`Group ${newStatus ? 'enabled' : 'disabled'}`, 'success')
          // 如果正在编辑此组，也更新编辑副本
          if (this.selectedGroup && this.selectedGroup.id === group.id) {
            this.editingGroup.enabled = newStatus
          }
        })
        .catch(error => {
          this.showMessage(`Failed to toggle group: ${error}`, 'error')
        })
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
        this.showMessage('Hosts applied successfully', 'success')
      } catch (error) {
        this.showMessage(`Failed to apply hosts: ${error}`, 'error')
      }
    },

    async refreshRemote() {
      try {
        await RefreshRemoteGroups()
        await this.loadHostGroups()
        this.showMessage('Remote groups refreshed', 'success')
      } catch (error) {
        this.showMessage(`Failed to refresh remote groups: ${error}`, 'error')
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
        // 获取系统host路径和内容
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
