<template>
  <div id="app">
    <div class="app-container">
      <!-- 左侧树形结构 -->
      <div class="sidebar">
        <div class="sidebar-header">
          <h2>Host Groups</h2>
          <button class="btn btn-add" @click="showAddGroupModal = true">+</button>
        </div>
        
        <div class="search-box">
          <input 
            v-model="searchQuery" 
            placeholder="Search groups..." 
            class="search-input"
          />
        </div>
        
        <div class="tree-view">
          <!-- 系统Host文件条目 -->
          <div 
            class="tree-item system-host-item"
            :class="{ 'active': selectedGroup && selectedGroup.id === 'system-host' }"
            @click="selectSystemHost()"
          >
            <div class="tree-item-content">
              <div class="item-icon">
                <span>⚙️</span>
              </div>
              <div class="item-details">
                <div class="item-name">System Host File</div>
                <div class="item-description">{{ systemHostPath }}</div>
              </div>
            </div>
          </div>
          
          <div 
            v-for="group in filteredGroups" 
            :key="group.id" 
            class="tree-item"
            :class="{ 'active': selectedGroup && selectedGroup.id === group.id }"
            @click="selectGroup(group)"
          >
            <div class="tree-item-content">
              <div class="item-icon">
                <span v-if="group.isRemote">🌐</span>
                <span v-else>📝</span>
              </div>
              <div class="item-details">
                <div class="item-name">{{ group.name }}</div>
                <div class="item-description">{{ group.description || 'No description' }}</div>
              </div>
              <div class="item-actions">
                <div 
                  class="switch-control" 
                  :class="{ 'enabled': group.enabled, 'disabled': !group.enabled }"
                  @click.stop="toggleGroupStatus(group)"
                  :title="group.enabled ? 'Click to disable' : 'Click to enable'"
                >
                  <div class="switch-slider">
                    <span class="switch-text">{{ group.enabled ? 'ON' : 'OFF' }}</span>
                  </div>
                </div>
                <button 
                  class="btn-icon" 
                  @click.stop="deleteGroup(group.id)"
                  title="Delete group"
                >
                  🗑️
                </button>
              </div>
            </div>
          </div>
          
          <div v-if="filteredGroups.length === 0 && !searchQuery" class="no-results">
            No host groups found
          </div>
        </div>
      </div>
      
      <!-- 右侧面板 - 编辑内容 -->
      <div class="main-panel">
        <!-- 远程Host预览 -->
        <div v-if="selectedGroup && selectedGroup.isRemote && selectedGroup.id !== 'system-host'" class="group-editor">
          <div class="editor-header">
            <h3>{{ selectedGroup.name }}</h3>
            <div class="header-actions">
              <button 
                class="btn btn-primary" 
                @click="refreshRemoteContent"
                :disabled="isRefreshingRemote"
              >
                {{ isRefreshingRemote ? 'Refreshing...' : 'Refresh Content' }}
              </button>
            </div>
          </div>
          
          <div class="editor-content">
            <div class="form-row">
              <div class="form-group-half">
                <label>Name *</label>
                <input 
                  v-model="editingGroup.name" 
                  @input="markAsDirty"
                  :readonly="selectedGroup.id.startsWith('system-')"
                />
              </div>
              <div class="form-group-half">
                <label>Description</label>
                <input 
                  v-model="editingGroup.description" 
                  @input="markAsDirty"
                  :readonly="selectedGroup.id.startsWith('system-')"
                />
              </div>
            </div>
            
            <div class="form-row">
              <div class="form-group-two-thirds">
                <label>URL</label>
                <input 
                  v-model="editingGroup.url" 
                  @input="markAsDirty"
                  placeholder="Remote hosts URL"
                  :readonly="selectedGroup.id.startsWith('system-')"
                />
              </div>
              <div class="form-group-one-third">
                <label>&nbsp;</label>
                <button 
                  class="btn btn-secondary btn-full-width" 
                  @click="fetchRemoteContent"
                  :disabled="!editingGroup.url || isFetchingRemote"
                >
                  {{ isFetchingRemote ? 'Getting...' : '获取host内容' }}
                </button>
              </div>
            </div>
            
            <div class="form-group" v-if="remoteContentPreview">
              <label>Remote Content (Preview)</label>
              <textarea 
                :value="remoteContentPreview" 
                readonly
                placeholder="Remote content will be displayed here after refresh..."
                rows="15"
                class="disabled-input"
              ></textarea>
            </div>
            
            <div class="group-meta">
              <p><strong>Created:</strong> {{ formatDate(selectedGroup.createdAt) }}</p>
              <p><strong>Last Updated:</strong> {{ formatDate(selectedGroup.updatedAt) }}</p>
              <p><strong>ID:</strong> {{ selectedGroup.id }}</p>
              <p><strong>Type:</strong> {{ selectedGroup.isRemote ? 'REMOTE' : 'LOCAL' }}</p>
              <p v-if="selectedGroup.lastUpdated"><strong>Last Fetched:</strong> {{ formatDate(selectedGroup.lastUpdated) }}</p>
            </div>
          </div>
        </div>
        
        <!-- 本地Host编辑 -->
        <div v-else-if="selectedGroup && !selectedGroup.isRemote && selectedGroup.id !== 'system-host'" class="group-editor">
          <div class="editor-header">
            <h3>{{ selectedGroup.name }}</h3>
            <div class="header-actions">
              <button 
                class="btn btn-primary" 
                @click="saveGroup"
                :disabled="!isDirty"
              >
                {{ isDirty ? 'Save Changes' : 'Saved' }}
              </button>
              <button 
                class="btn btn-secondary" 
                @click="cancelEdit"
                v-if="isDirty"
              >
                Cancel
              </button>
            </div>
          </div>
          
          <div class="editor-content">
            <div class="form-group">
              <label>Name *</label>
              <input 
                v-model="editingGroup.name" 
                @input="markAsDirty"
              />
            </div>
            
            <div class="form-group">
              <label>Description</label>
              <input 
                v-model="editingGroup.description" 
                @input="markAsDirty"
              />
            </div>
            
            <div v-if="editingGroup.isRemote" class="form-group">
              <label>URL</label>
              <input 
                v-model="editingGroup.url" 
                @input="markAsDirty"
                placeholder="Remote hosts URL"
              />
            </div>
            
            <div v-if="!editingGroup.isRemote" class="form-group">
              <label>Content</label>
              <textarea 
                v-model="editingGroup.content" 
                @input="markAsDirty"
                placeholder="Enter host entries here..."
                rows="20"
              ></textarea>
            </div>
            
            <div class="group-meta">
              <p><strong>Created:</strong> {{ formatDate(selectedGroup.createdAt) }}</p>
              <p><strong>Last Updated:</strong> {{ formatDate(selectedGroup.updatedAt) }}</p>
              <p><strong>ID:</strong> {{ selectedGroup.id }}</p>
              <p><strong>Type:</strong> {{ selectedGroup.isRemote ? 'REMOTE' : 'LOCAL' }}</p>
            </div>
          </div>
        </div>
        
        <!-- 系统Host文件预览 -->
        <div v-else-if="selectedGroup && selectedGroup.id === 'system-host'" class="group-editor">
          <div class="editor-header">
            <h3>System Host File</h3>
            <div class="header-actions">
              <button class="btn btn-secondary" @click="refreshSystemHost">Refresh</button>
            </div>
          </div>
          
          <div class="editor-content">
            <div class="form-group">
              <label>File Path</label>
              <input 
                :value="systemHostPath" 
                readonly
                class="disabled-input"
              />
            </div>
            
            <div class="form-group">
              <label>Content (Read-Only)</label>
              <textarea 
                :value="systemHostContent" 
                readonly
                placeholder="System host file content will be displayed here..."
                rows="20"
                class="disabled-input"
              ></textarea>
            </div>
            
            <div class="group-meta">
              <p><strong>Type:</strong> System File</p>
              <p><strong>Editable:</strong> No (Requires admin privileges)</p>
            </div>
          </div>
        </div>
        
        <div v-else class="welcome-panel">
          <h3>Select a host group to edit</h3>
          <p>Choose a host group from the left sidebar to view and edit its content.</p>
        </div>
      </div>
    </div>
    
    <!-- 顶部操作栏 -->
    <div class="top-bar">
      <button class="btn btn-primary" @click="applyHosts" title="Apply all enabled host groups to system">
        Apply Hosts to System
      </button>
      <button class="btn btn-secondary" @click="refreshRemote" title="Refresh all remote host groups">
        Refresh Remote Groups
      </button>
      <button class="btn btn-info" @click="loadHostGroups" title="Refresh the list">
        Refresh
      </button>
    </div>
    
    <!-- 添加组模态框 -->
    <div v-if="showAddGroupModal" class="modal-overlay" @click="closeAddGroupModal">
      <div class="modal-content" @click.stop>
        <h3>Add New Host Group</h3>
        
        <!-- 选择Host类型 -->
        <div v-if="!newGroup.typeSelected" class="form-group">
          <label>Select Host Type</label>
          <div class="host-type-selection">
            <button class="btn btn-option" @click="selectHostType(false)">Local Host</button>
            <button class="btn btn-option" @click="selectHostType(true)">Remote Host</button>
          </div>
        </div>
        
        <!-- 填写Host信息 -->
        <div v-if="newGroup.typeSelected">
          <div class="form-group">
            <label>Name *</label>
            <input 
              v-model="newGroup.name" 
              placeholder="Display name"
            />
          </div>
          <div class="form-group">
            <label>Description</label>
            <input 
              v-model="newGroup.description" 
              placeholder="Description"
            />
          </div>
          
          <!-- 如果是远程Host，显示URL字段 -->
          <div v-if="newGroup.isRemote" class="form-group">
            <label>URL</label>
            <input 
              v-model="newGroup.url" 
              placeholder="Remote hosts URL"
            />
          </div>
          
          <div class="form-row">
            <button class="btn btn-secondary" @click="resetHostTypeSelection">Back</button>
            <button class="btn btn-primary" @click="addGroup">Add Group</button>
          </div>
        </div>
        
        <div v-if="newGroup.typeSelected" class="modal-actions">
          <button class="btn btn-secondary" @click="closeAddGroupModal">Cancel</button>
        </div>
      </div>
    </div>
    
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
  GetRemoteContent
} from '../wailsjs/go/main/App'

export default {
  name: 'App',
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
  computed: {
    filteredGroups() {
      if (!this.searchQuery) {
        return this.groups
      }
      const query = this.searchQuery.toLowerCase()
      return this.groups.filter(group => 
        (group.name && group.name.toLowerCase().includes(query)) || 
        (group.description && group.description.toLowerCase().includes(query)) ||
        (group.id && group.id.toLowerCase().includes(query))
      )
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

    async addGroup() {
      if (!this.newGroup.name) {
        this.showMessage('Name is required', 'error')
        return
      }

      try {
        await AddHostGroup({ ...this.newGroup })
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

    async saveGroup() {
      try {
        await UpdateHostGroup(this.editingGroup)
        // 更新主列表中的组
        const index = this.groups.findIndex(g => g.id === this.editingGroup.id)
        if (index !== -1) {
          this.groups[index] = { ...this.editingGroup }
          this.selectedGroup = { ...this.editingGroup }
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

    selectHostType(isRemote) {
      this.newGroup.isRemote = isRemote;
      this.newGroup.typeSelected = true;
    },

    resetHostTypeSelection() {
      this.newGroup.typeSelected = false;
      this.newGroup.name = '';
      this.newGroup.description = '';
      this.newGroup.url = '';
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
        // 从指定URL获取远程内容
        this.remoteContentPreview = await GetRemoteContent(this.selectedGroup.url)
        
        // 更新本地组的content字段
        this.selectedGroup.content = this.remoteContentPreview
        this.editingGroup.content = this.remoteContentPreview
        
        this.showMessage('Remote content fetched successfully', 'success')
      } catch (error) {
        this.showMessage(`Failed to fetch remote content: ${error}`, 'error')
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
        this.showMessage('Remote content fetched successfully', 'success');
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

    formatDate(dateString) {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleString()
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

.sidebar {
  width: 350px;
  background: #f8f9fa;
  border-right: 1px solid #dee2e6;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-header h2 {
  margin: 0;
  color: #495057;
}

.btn-add {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  border: none;
  background: #007bff;
  color: white;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-box {
  padding: 10px 15px;
  border-bottom: 1px solid #dee2e6;
}

.search-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.tree-view {
  flex: 1;
  overflow-y: auto;
  padding: 10px 0;
}

.tree-item {
  padding: 12px 15px;
  cursor: pointer;
  border-left: 3px solid transparent;
  transition: all 0.2s;
}

.tree-item:hover {
  background: #e9ecef;
}

.tree-item.active {
  background: #e3f2fd;
  border-left-color: #2196f3;
}

.tree-item-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.item-icon {
  font-size: 18px;
  min-width: 24px;
}

.item-details {
  flex: 1;
  overflow: hidden;
}

.item-name {
  font-weight: 600;
  color: #212529;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-description {
  font-size: 12px;
  color: #6c757d;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-indicator {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  cursor: pointer;
}

.status-indicator.enabled {
  background: #28a745;
}

.status-indicator.disabled {
  background: #dc3545;
}

.btn-icon {
  background: none;
  border: none;
  font-size: 14px;
  cursor: pointer;
  padding: 4px;
  border-radius: 3px;
}

.btn-icon:hover {
  background: rgba(0, 0, 0, 0.1);
}

.no-results {
  padding: 20px;
  text-align: center;
  color: #6c757d;
  font-style: italic;
}

.main-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: white;
}

.welcome-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: #6c757d;
}

.welcome-panel h3 {
  margin-bottom: 10px;
  color: #495057;
}

.group-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.editor-header {
  padding: 20px;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.editor-header h3 {
  margin: 0;
  color: #495057;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.editor-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 600;
  color: #495057;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #80bdff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

.form-group input:disabled,
.form-group input.disabled-input {
  background-color: #e9ecef;
  cursor: not-allowed;
}

.form-group textarea {
  font-family: monospace;
  resize: vertical;
  min-height: 150px;
}

.group-meta {
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #eee;
  font-size: 14px;
}

.group-meta p {
  margin: 5px 0;
}

.enabled {
  color: #28a745;
  font-weight: 600;
}

.disabled {
  color: #dc3545;
  font-weight: 600;
}

.top-bar {
  padding: 10px 20px;
  background: #f8f9fa;
  border-top: 1px solid #dee2e6;
  display: flex;
  gap: 10px;
  justify-content: flex-start;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #495057;
}

.modal-actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
  text-decoration: none;
  display: inline-block;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #0056b3;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #545b62;
}

.btn-info {
  background-color: #17a2b8;
  color: white;
}

.btn-info:hover {
  background-color: #117a8b;
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

/* Switch 控件样式 */
.switch-control {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 24px;
}

.switch-control input {
  opacity: 0;
  width: 0;
  height: 0;
}

.switch-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
  border-radius: 24px;
}

.switch-slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

.switch-control.enabled .switch-slider {
  background-color: #2196F3;
}

.switch-control.enabled .switch-slider:before {
  transform: translateX(26px);
}

.switch-control input:checked + .switch-slider {
  background-color: #2196F3;
}

.switch-control input:checked + .switch-slider:before {
  transform: translateX(26px);
}

.switch-text {
  display: none;
}

/* 系统Host文件条目样式 */
.system-host-item {
  background-color: #e3f2fd;
  border-left: 3px solid #2196f3;
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

/* 表单行布局 */
.form-row {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.form-group-half {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.form-group-two-thirds {
  flex: 2;
  display: flex;
  flex-direction: column;
}

.form-group-one-third {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.btn-full-width {
  width: 100%;
  height: 40px;
  margin-top: 19px; /* Align with input field */
}

.host-type-selection {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.host-type-selection .btn {
  flex: 1;
}

.btn-option {
  padding: 12px 16px;
  border: 2px solid #007bff;
  background-color: white;
  color: #007bff;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.btn-option:hover {
  background-color: #007bff;
  color: white;
}

.btn-option.selected {
  background-color: #007bff;
  color: white;
}
</style>
