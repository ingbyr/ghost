<template>
  <div class="group-editor">
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
</template>

<script>
export default {
  name: 'SystemHostPreview',
  props: {
    systemHostPath: {
      type: String,
      default: ''
    },
    systemHostContent: {
      type: String,
      default: ''
    }
  },
  emits: ['refresh-system-host'],
  methods: {
    refreshSystemHost() {
      this.$emit('refresh-system-host');
    }
  }
}
</script>

<style scoped>
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
</style>