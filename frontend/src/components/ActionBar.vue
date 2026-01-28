<template>
  <div class="top-bar">
    <button class="btn btn-secondary" @click="refreshRemote" title="Refresh all remote host groups">
      Refresh Remote Groups
    </button>
    <button class="btn btn-info" @click="backupNow" title="Create manual backup">
      Backup Now
    </button>
  </div>
</template>

<script>
export default {
  name: 'ActionBar',
  emits: ['refresh-remote', 'refresh-list', 'backup-now'],
  methods: {
    refreshRemote() {
      this.$emit('refresh-remote');
    },
    refreshList() {
      this.$emit('refresh-list');
    },
    backupNow() {
      this.$emit('backup-now');
    }
  }
}
</script>

<style scoped>
.top-bar {
  padding: 10px 20px;
  background: #f8f9fa;
  border-top: 1px solid #dee2e6;
  display: flex;
  gap: 10px;
  justify-content: flex-start;
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
</style>