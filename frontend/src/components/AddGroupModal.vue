<template>
  <div v-if="showModal" class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <h3>{{ t('components.addGroupModal.title') }}</h3>
      
      <!-- 选择Host类型 -->
      <div v-if="!localNewGroup.typeSelected" class="form-group">
        <label>{{ t('components.addGroupModal.selectHostType') }}</label>
        <div class="host-type-selection">
          <button class="btn btn-option" @click="selectHostType(false)">{{ t('components.addGroupModal.localHost') }}</button>
          <button class="btn btn-option" @click="selectHostType(true)">{{ t('components.addGroupModal.remoteHost') }}</button>
        </div>
      </div>
      
      <!-- 填写Host信息 -->
      <div v-if="localNewGroup.typeSelected">
        <el-row :gutter="20">
          <el-col :span="12">
            <div class="form-group">
              <label>{{ t('common.name') }} *</label>
              <input 
                v-model="localNewGroup.name" 
                :placeholder="t('common.name')"
              />
            </div>
          </el-col>
          <el-col :span="12">
            <div class="form-group">
              <label>{{ t('common.description') }}</label>
              <input 
                v-model="localNewGroup.description" 
                :placeholder="t('common.description')"
              />
            </div>
          </el-col>
        </el-row>
        
        <!-- 如果是远程Host，显示URL字段 -->
        <div v-if="localNewGroup.isRemote" class="form-group">
          <label>{{ t('common.url') }}</label>
          <input 
            v-model="localNewGroup.url" 
            :placeholder="t('components.remoteHostEditor.remoteContentPlaceholder')"
          />
        </div>
        
        <div class="form-row">
          <button class="btn btn-secondary" @click="resetHostTypeSelection">{{ t('common.back') }}</button>
          <button class="btn btn-primary" @click="addGroup">{{ t('components.addGroupModal.addGroup') }}</button>
        </div>
      </div>
      
      <div v-if="localNewGroup.typeSelected" class="modal-actions">
        <button class="btn btn-secondary" @click="closeModal">{{ t('components.addGroupModal.cancel') }}</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ElRow, ElCol } from 'element-plus';
import { useI18n } from 'vue-i18n';

export default {
  name: 'AddGroupModal',
  components: {
    ElRow,
    ElCol
  },
  props: {
    showModal: {
      type: Boolean,
      default: false
    },
    newGroup: {
      type: Object,
      default: () => ({
        name: '',
        description: '',
        content: '',
        enabled: false,
        isRemote: false,
        url: ''
      })
    }
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  emits: ['close-modal', 'add-group', 'update:newGroup'],
  data() {
    return {
      localNewGroup: { ...this.newGroup }
    }
  },
  watch: {
    newGroup: {
      handler(newVal) {
        this.localNewGroup = { ...newVal };
      },
      deep: true
    }
  },
  methods: {
    closeModal() {
      this.$emit('close-modal');
    },
    addGroup() {
      this.$emit('add-group', this.localNewGroup);
    },
    selectHostType(isRemote) {
      this.localNewGroup.isRemote = isRemote;
      this.localNewGroup.typeSelected = true;
      this.$emit('update:newGroup', this.localNewGroup);
    },
    resetHostTypeSelection() {
      this.localNewGroup.typeSelected = false;
      this.localNewGroup.name = '';
      this.localNewGroup.description = '';
      this.localNewGroup.url = '';
      this.$emit('update:newGroup', this.localNewGroup);
    }
  }
}
</script>

<style scoped>
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

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #80bdff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
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

.modal-actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  justify-content: flex-end;
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

.form-row {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}
</style>