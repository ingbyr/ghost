<template>
  <div class="main-panel">
    <RemoteHostEditor 
      v-if="selectedGroup && selectedGroup.isRemote && selectedGroup.id !== 'system-host'"
      :group="selectedGroup"
      :editing-group="editingGroup"
      :remote-content-preview="remoteContentPreview"
      :is-fetching-remote="isFetchingRemote"
      @save-group="$emit('save-group', $event)"
      @cancel-edit="$emit('cancel-edit')"
      @fetch-remote-content="$emit('fetch-remote-content')"
      @mark-as-dirty="$emit('mark-as-dirty')"
      @apply-hosts="$emit('apply-hosts')"
    />
    
    <LocalHostEditor
      v-else-if="selectedGroup && !selectedGroup.isRemote && selectedGroup.id !== 'system-host'"
      :group="selectedGroup"
      :editing-group="editingGroup"
      @save-group="$emit('save-group', $event)"
      @cancel-edit="$emit('cancel-edit')"
      @mark-as-dirty="$emit('mark-as-dirty')"
    />
    
    <SystemHostPreview
      v-else-if="selectedGroup && selectedGroup.id === 'system-host'"
      :system-host-path="systemHostPath"
      :system-host-content="systemHostContent"
      @refresh-system-host="$emit('refresh-system-host')"
    />
    
    <div v-else class="welcome-panel">
      <h3>Select a host group to edit</h3>
      <p>Choose a host group from the left sidebar to view and edit its content.</p>
    </div>
  </div>
</template>

<script>
import RemoteHostEditor from './RemoteHostEditor.vue';
import LocalHostEditor from './LocalHostEditor.vue';
import SystemHostPreview from './SystemHostPreview.vue';

export default {
  name: 'MainPanel',
  components: {
    RemoteHostEditor,
    LocalHostEditor,
    SystemHostPreview
  },
  props: {
    selectedGroup: {
      type: Object,
      default: null
    },
    editingGroup: {
      type: Object,
      default: () => ({})
    },
    systemHostPath: {
      type: String,
      default: ''
    },
    systemHostContent: {
      type: String,
      default: ''
    },
    remoteContentPreview: {
      type: String,
      default: ''
    },
    isFetchingRemote: {
      type: Boolean,
      default: false
    }
  },
  emits: [
    'save-group', 
    'cancel-edit', 
    'fetch-remote-content', 
    'mark-as-dirty',
    'refresh-system-host',
    'apply-hosts'
  ]
}
</script>

<style scoped>
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
</style>