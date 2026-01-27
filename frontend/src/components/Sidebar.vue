<template>
  <div class="sidebar">
    <div class="sidebar-header">
      <h2>Host Groups</h2>
      <button class="btn btn-add" @click="$emit('open-add-modal')">+</button>
    </div>
    
    <div class="search-box">
      <input 
        :value="searchQuery"
        @input="$emit('update:search-query', $event.target.value)"
        placeholder="Search groups..." 
        class="search-input"
      />
    </div>
    
    <div class="tree-view">
      <!-- 系统Host文件条目 -->
      <div 
        class="tree-item system-host-item"
        :class="{ 'active': selectedGroup && selectedGroup.id === 'system-host' }"
        @click="selectSystemHost"
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
        @click="$emit('select-group', group)"
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
              @click.stop="$emit('toggle-status', group)"
              :title="group.enabled ? 'Click to disable' : 'Click to enable'"
            >
              <div class="switch-slider">
                <span class="switch-text">{{ group.enabled ? 'ON' : 'OFF' }}</span>
              </div>
            </div>
            <button 
              class="btn-icon" 
              @click.stop="$emit('delete-group', group.id)"
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
</template>

<script>
export default {
  name: 'Sidebar',
  props: {
    groups: {
      type: Array,
      required: true
    },
    selectedGroup: {
      type: Object,
      default: null
    },
    systemHostPath: {
      type: String,
      default: ''
    },
    searchQuery: {
      type: String,
      default: ''
    }
  },
  emits: ['select-group', 'toggle-status', 'delete-group', 'open-add-modal', 'update:search-query'],
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
  methods: {
    selectSystemHost() {
      this.$emit('select-system-host');
    }
  }
}
</script>

<style scoped>
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
</style>