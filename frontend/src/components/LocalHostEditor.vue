<template>
  <div class="group-editor">
    <div class="editor-header">
      <h3>{{ group.name }}</h3>
    </div>
    
    <div class="editor-content">
      <el-row :gutter="20" class="form-row">
        <el-col :span="12">
          <div class="form-group">
            <label>{{ t('common.name') }} *</label>
            <input 
              v-model="localEditingGroup.name" 
              @input="markAsDirty"
              :disabled="isReadOnly"
            />
          </div>
        </el-col>
        <el-col :span="12">
          <div class="form-group">
            <label>{{ t('common.description') }}</label>
            <input 
              v-model="localEditingGroup.description" 
              @input="markAsDirty"
              :disabled="isReadOnly"
            />
          </div>
        </el-col>
      </el-row>
      
      <div v-if="localEditingGroup.isRemote" class="form-group">
        <label>{{ t('common.url') }}</label>
        <input 
          v-model="localEditingGroup.url" 
          @input="markAsDirty"
          :placeholder="t('components.remoteHostEditor.remoteContentPlaceholder')"
          :disabled="isReadOnly"
        />
      </div>
      
      <div v-if="!localEditingGroup.isRemote" class="form-group">
        <label>
          {{ t('common.content') }}
          <span v-if="isReadOnly" class="read-only-hint"> {{ t('components.localhostEditor.readOnlyHint') }}</span>
        </label>
        <!-- 启用时为只读状态，禁用时为可编辑状态 -->
        <textarea 
          :value="localEditingGroup.content"
          @input="!isReadOnly ? handleContentInput($event) : null"
          @blur="!isReadOnly ? handleContentBlur($event) : null"
          :readonly="isReadOnly"
          :placeholder="t('components.localhostEditor.enterHostEntries')"
          rows="20"
        ></textarea>
      </div>
      
      <div class="group-meta">
        <p><strong>{{ t('common.created') }}:</strong> {{ formatDate(group.createdAt) }}</p>
        <p><strong>{{ t('common.lastUpdated') }}:</strong> {{ formatDate(group.updatedAt) }}</p>
        <p><strong>{{ t('common.id') }}:</strong> {{ group.id }}</p>
        <p><strong>{{ t('common.type') }}:</strong> {{ group.isRemote ? t('common.remote') : t('common.local') }}</p>
        <p><strong>{{ t('common.status') }}:</strong> 
          <span :class="{'status-enabled': group.enabled, 'status-disabled': !group.enabled}">
            {{ group.enabled ? t('status.enabled') : t('status.disabled') }}
          </span>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { ElRow, ElCol } from 'element-plus';
import { useI18n } from 'vue-i18n';

export default {
  name: 'LocalHostEditor',
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
    }
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  emits: ['save-group', 'cancel-edit', 'mark-as-dirty'],
  data() {
    return {
      localEditingGroup: { ...this.editingGroup }
    }
  },
  computed: {
    isDirty() {
      return JSON.stringify(this.localEditingGroup) !== JSON.stringify(this.group);
    },
    // 计算属性：如果组已启用，则不允许编辑（仅适用于本地组）
    isReadOnly() {
      return this.group.enabled === true && !this.group.isRemote;
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
    handleContentInput(event) {
      // 更新本地编辑组的内容
      this.localEditingGroup.content = event.target.value;
      this.markAsDirty();
    },
    handleContentBlur(event) {
      // 当内容编辑区域失去焦点时自动保存，但只在不是只读模式时
      if (!this.isReadOnly && this.isDirty) {
        this.saveGroup();
      }
    },
    saveGroup() {
      // 如果是只读模式，则不允许保存
      if (this.isReadOnly) {
        return;
      }
      this.$emit('save-group', this.localEditingGroup);
    },
    cancelEdit() {
      this.$emit('cancel-edit');
    },
    markAsDirty() {
      // 如果是只读模式，则不允许标记为dirty
      if (this.isReadOnly) {
        return;
      }
      this.$emit('mark-as-dirty');
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

.read-only-hint {
  color: #6c757d;
  font-size: 0.9em;
  font-weight: normal;
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

.disabled-input {
  background-color: #e9ecef;
  cursor: not-allowed;
}

.floating-save-btn {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  gap: 10px;
  z-index: 1000;
  background: white;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border: 1px solid #dee2e6;
}
</style>