<template>
  <div class="top-bar">
    <button class="btn btn-secondary" @click="refreshRemote" :title="t('components.actionBar.refreshRemoteGroupsTooltip')">
      {{ t('components.actionBar.refreshRemoteGroups') }}
    </button>
    <button class="btn btn-info" @click="backupNow" :title="t('components.actionBar.backupNowTooltip')">
      {{ t('components.actionBar.backupNow') }}
    </button>
    <button class="btn btn-warning" @click="restoreBackup" :title="t('components.actionBar.restoreBackupTooltip')">
      {{ t('components.actionBar.restoreBackup') }}
    </button>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';

export default {
  name: 'ActionBar',
  setup() {
    const { t } = useI18n();
    return { t };
  },
  emits: ['refresh-remote', 'refresh-list', 'backup-now', 'restore-backup'],
  methods: {
    refreshRemote() {
      this.$emit('refresh-remote');
    },
    refreshList() {
      this.$emit('refresh-list');
    },
    backupNow() {
      this.$emit('backup-now');
    },
    restoreBackup() {
      this.$emit('restore-backup');
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

.btn-warning {
  background-color: #ffc107;
  color: #212529;
}

.btn-warning:hover {
  background-color: #e0a800;
}
</style>