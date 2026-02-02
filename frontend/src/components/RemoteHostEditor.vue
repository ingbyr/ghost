<template>
  <div class="group-editor">
    <div class="editor-header">
      <h3>{{ group.name }}</h3>
    </div>
    
    <div class="editor-content">
      <el-row :gutter="20" class="form-row">
        <el-col :span="8">
          <div class="form-group">
            <label>{{ t('common.name') }} *</label>
            <input 
              v-model="localEditingGroup.name" 
              @input="markAsDirty"
              @blur="autoSave"
              :readonly="group.id.startsWith('system-')"
            />
          </div>
        </el-col>
        <el-col :span="8">
          <div class="form-group">
            <label>{{ t('common.description') }}</label>
            <input 
              v-model="localEditingGroup.description" 
              @input="markAsDirty"
              @blur="autoSave"
              :readonly="group.id.startsWith('system-')"
            />
          </div>
        </el-col>
        <el-col :span="8">
          <div class="form-group">
            <label>{{ t('components.remoteHostEditor.refreshInterval') }}</label>
            <el-select 
              v-model="localEditingGroup.refreshInterval" 
              @change="autoSave"
              :placeholder="t('components.remoteHostEditor.selectRefreshInterval')"
              class="full-width"
              :disabled="group.id.startsWith('system-')"
            >
              <el-option 
                :label="t('components.remoteHostEditor.refreshDisabled')" 
                :value="0" 
              />
              <el-option 
                :label="t('components.remoteHostEditor.everyHour')" 
                :value="3600" 
              />
              <el-option 
                :label="t('components.remoteHostEditor.every4Hours')" 
                :value="14400" 
              />
              <el-option 
                :label="t('components.remoteHostEditor.every8Hours')" 
                :value="28800" 
              />
              <el-option 
                :label="t('components.remoteHostEditor.every24Hours')" 
                :value="86400" 
              />
            </el-select>
          </div>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" class="form-row">
        <el-col :span="16">
          <div class="form-group">
            <label>{{ t('common.url') }}</label>
            <input 
              v-model="localEditingGroup.url" 
              @input="markAsDirty"
              @blur="autoSave"
              :placeholder="t('components.remoteHostEditor.remoteContentPlaceholder')"
              :readonly="group.id.startsWith('system-')"
            />
          </div>
        </el-col>
        <el-col :span="8">
          <div class="form-group">
            <button
              class="btn btn-secondary btn-full-width" 
              @click="fetchRemoteContent"
              :disabled="!localEditingGroup.url || isFetchingRemote"
            >
              {{ isFetchingRemote ? t('components.remoteHostEditor.getting') : t('components.remoteHostEditor.fetchHostContent') }}
            </button>
          </div>
        </el-col>
      </el-row>
      
      <div class="form-group" v-if="remoteContentPreview">
        <label>{{ t('components.remoteHostEditor.remoteContentPreview') }}</label>
        <textarea 
          :value="remoteContentPreview" 
          readonly
          :placeholder="t('components.remoteHostEditor.remoteContentPlaceholder')"
          rows="15"
          class="disabled-input"
        ></textarea>
      </div>
      
      <div class="group-meta">
        <p><strong>{{ t('common.created') }}:</strong> {{ formatDate(group.createdAt) }}</p>
        <p><strong>{{ t('common.lastUpdated') }}:</strong> {{ formatDate(group.updatedAt) }}</p>
        <p><strong>{{ t('common.id') }}:</strong> {{ group.id }}</p>
        <p><strong>{{ t('common.type') }}:</strong> {{ group.isRemote ? t('common.remote') : t('common.local') }}</p>
        <p v-if="group.lastUpdated"><strong>{{ t('messages.lastFetched') }}:</strong> {{ formatDate(group.lastUpdated) }}</p>
      </div>
    </div>
    

  </div>
</template>

<script>
import { ElRow, ElCol } from 'element-plus';
import { useI18n } from 'vue-i18n';

export default {
  name: 'RemoteHostEditor',
  components: {
    ElRow,
    ElCol
  },
  props: {
    group: {
      type: Object,
      required: true
    },
    editingGroup: {
      type: Object,
      required: true
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
  setup() {
    const { t } = useI18n();
    return { t };
  },
  emits: ['save-group', 'cancel-edit', 'fetch-remote-content', 'mark-as-dirty', 'apply-hosts'],
  data() {
    return {
      localEditingGroup: { ...this.editingGroup }
    }
  },
  computed: {
    isDirty() {
      return JSON.stringify(this.localEditingGroup) !== JSON.stringify(this.group);
    }
  },
  watch: {
    editingGroup: {
      handler(newVal) {
        this.localEditingGroup = { ...newVal };
      },
      deep: true
    }
  },
  methods: {
    saveGroup() {
      this.$emit('save-group', this.localEditingGroup);
      // 保存后如果组已启用则应用Hosts
      if (this.localEditingGroup.enabled) {
        this.$emit('apply-hosts');
      }
    },
    cancelEdit() {
      this.$emit('cancel-edit');
    },
    fetchRemoteContent() {
      this.$emit('fetch-remote-content');
      // 获取远程内容后需要应用Hosts，如果该组已启用
      this.$emit('apply-hosts');
    },
    markAsDirty() {
      this.$emit('mark-as-dirty');
    },
    async autoSave() {
      // 检查是否有更改
      if (this.isDirty) {
        await this.saveGroup();
      }
    },
    formatDate(dateString) {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleString()
    }
  }
}
</script>

<style scoped>
.group-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  position: relative;
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
  text-align: left; /* 左对齐标签文本 */
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

.btn-full-width {
  width: 100%;
  height: 40px;
  margin-top: 28px; /* Adjusted to align with input field */
}

.full-width {
  width: 100%;
}

.disabled-input {
  background-color: #e9ecef;
  cursor: not-allowed;
}


</style>