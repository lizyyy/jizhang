// API 基础URL
const API_BASE = '/api';

// 类型映射
const TYPE_MAP = {
    expense: { label: '支出', icon: '💸', color: '#ff6b6b' },
    income: { label: '收入', icon: '💰', color: '#51cf66' },
    transfer: { label: '转账', icon: '💱', color: '#339af0' },
    loan: { label: '贷款', icon: '💳', color: '#fd7e14' }
};

// 全局数据
let categories = [];
let records = [];

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    initDateInputs();
    loadCategories();
    loadRecords();
    loadStatistics();
    bindEvents();
    bindTabEvents();
});

// 初始化日期输入框
function initDateInputs() {
    const today = new Date().toISOString().split('T')[0];
    document.getElementById('date').value = today;
}

// 绑定标签切换事件
function bindTabEvents() {
    const tabs = document.querySelectorAll('.nav-tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const targetTab = tab.dataset.tab;
            
            // 更新标签样式
            tabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');
            
            // 切换内容
            document.querySelectorAll('.tab-content').forEach(content => {
                content.classList.remove('active');
            });
            document.getElementById(targetTab + 'Tab').classList.add('active');
            
            // 如果切换到分类管理页面，渲染分类
            if (targetTab === 'categories') {
                renderCategoryManagement();
            }
        });
    });
}

// 绑定事件
function bindEvents() {
    // 表单提交
    document.getElementById('recordForm').addEventListener('submit', handleAddRecord);
    
    // 重置按钮
    document.getElementById('resetBtn').addEventListener('click', () => {
        document.getElementById('recordForm').reset();
        initDateInputs();
    });

    // 类型改变时更新分类
    document.getElementById('type').addEventListener('change', updateCategoryOptions);
    document.getElementById('editType').addEventListener('change', updateEditCategoryOptions);

    // 分类改变时更新子分类
    document.getElementById('category').addEventListener('change', updateSubcategoryOptions);
    document.getElementById('editCategory').addEventListener('change', updateEditSubcategoryOptions);

    // 筛选按钮
    document.getElementById('filterBtn').addEventListener('click', handleFilter);
    document.getElementById('clearFilterBtn').addEventListener('click', clearFilter);

    // 编辑表单
    document.getElementById('editForm').addEventListener('submit', handleUpdateRecord);
    document.getElementById('closeModal').addEventListener('click', closeEditModal);
    document.getElementById('cancelEdit').addEventListener('click', closeEditModal);

    // 点击模态框外部关闭
    document.getElementById('editModal').addEventListener('click', (e) => {
        if (e.target.id === 'editModal') {
            closeEditModal();
        }
    });

    // 分类管理表单
    document.getElementById('editCategoryForm').addEventListener('submit', handleUpdateCategory);
    document.getElementById('editSubcategoryForm').addEventListener('submit', handleUpdateSubcategory);
    document.getElementById('addSubcategoryForm').addEventListener('submit', handleAddSubcategory);
}

// 加载分类
async function loadCategories() {
    try {
        const response = await fetch(`${API_BASE}/categories`);
        const result = await response.json();
        
        if (result.Errno === 0) {
            categories = result.data;
            updateCategoryOptions();
            updateFilterCategoryOptions();
        } else {
            showError('加载分类失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('加载分类失败: ' + error.message);
    }
}

// 更新分类选项
function updateCategoryOptions() {
    const type = document.getElementById('type').value;
    const categorySelect = document.getElementById('category');
    
    const filteredCategories = categories.filter(cat => cat.type === type);
    
    categorySelect.innerHTML = '<option value="">请选择分类</option>';
    filteredCategories.forEach(cat => {
        const option = document.createElement('option');
        option.value = cat.id;
        option.textContent = cat.name;
        categorySelect.appendChild(option);
    });
    
    // 清空子分类
    document.getElementById('subcategory').innerHTML = '<option value="">请选择子分类</option>';
}

// 更新编辑表单分类选项
function updateEditCategoryOptions() {
    const type = document.getElementById('editType').value;
    const categorySelect = document.getElementById('editCategory');
    const currentValue = categorySelect.value;
    
    const filteredCategories = categories.filter(cat => cat.type === type);
    
    categorySelect.innerHTML = '<option value="">请选择分类</option>';
    filteredCategories.forEach(cat => {
        const option = document.createElement('option');
        option.value = cat.id;
        option.textContent = cat.name;
        categorySelect.appendChild(option);
    });
    
    // 恢复之前的选择（如果还在列表中）
    if (currentValue) {
        categorySelect.value = currentValue;
    }
    
    updateEditSubcategoryOptions();
}

// 更新子分类选项
function updateSubcategoryOptions() {
    const categoryId = document.getElementById('category').value;
    const subcategorySelect = document.getElementById('subcategory');
    
    subcategorySelect.innerHTML = '<option value="">请选择子分类</option>';
    
    if (!categoryId) return;
    
    const category = categories.find(cat => cat.id == categoryId);
    if (category && category.subcategories) {
        category.subcategories.forEach(sub => {
            const option = document.createElement('option');
            option.value = sub.id;
            option.textContent = sub.name;
            subcategorySelect.appendChild(option);
        });
    }
}

// 更新编辑表单子分类选项
function updateEditSubcategoryOptions() {
    const categoryId = document.getElementById('editCategory').value;
    const subcategorySelect = document.getElementById('editSubcategory');
    const currentValue = subcategorySelect.value;
    
    subcategorySelect.innerHTML = '<option value="">请选择子分类</option>';
    
    if (!categoryId) return;
    
    const category = categories.find(cat => cat.id == categoryId);
    if (category && category.subcategories) {
        category.subcategories.forEach(sub => {
            const option = document.createElement('option');
            option.value = sub.id;
            option.textContent = sub.name;
            subcategorySelect.appendChild(option);
        });
    }
    
    // 恢复之前的选择（如果还在列表中）
    if (currentValue) {
        subcategorySelect.value = currentValue;
    }
}

// 更新筛选分类选项
function updateFilterCategoryOptions() {
    const categorySelect = document.getElementById('filterCategory');
    
    categorySelect.innerHTML = '<option value="">全部</option>';
    categories.forEach(cat => {
        const option = document.createElement('option');
        option.value = cat.id;
        option.textContent = `${TYPE_MAP[cat.type].label} - ${cat.name}`;
        categorySelect.appendChild(option);
    });
}

// 加载记录
async function loadRecords(filters = {}) {
    try {
        let url = `${API_BASE}/records`;
        const params = new URLSearchParams();
        
        if (filters.type) params.append('type', filters.type);
        if (filters.categoryId) params.append('categoryId', filters.categoryId);
        if (filters.startDate) params.append('startDate', filters.startDate);
        if (filters.endDate) params.append('endDate', filters.endDate);
        if (filters.keyword) params.append('keyword', filters.keyword);
        
        if (params.toString()) {
            url += '?' + params.toString();
        }
        
        const response = await fetch(url);
        const result = await response.json();
        
        if (result.Errno === 0) {
            records = result.data;
            renderRecords();
        } else {
            showError('加载记录失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('加载记录失败: ' + error.message);
    }
}

// 加载统计
async function loadStatistics() {
    try {
        const response = await fetch(`${API_BASE}/records/statistics`);
        const result = await response.json();
        
        if (result.Errno === 0) {
            const stats = result.data;
            document.getElementById('totalExpense').textContent = formatMoney(stats.expense || 0);
            document.getElementById('totalIncome').textContent = formatMoney(stats.income || 0);
            const balance = (stats.income || 0) - (stats.expense || 0);
            document.getElementById('balance').textContent = formatMoney(balance);
        }
    } catch (error) {
        console.error('加载统计失败:', error);
    }
}

// 渲染记录列表
function renderRecords() {
    const container = document.getElementById('recordsList');
    
    if (records.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">📭</div>
                <p>暂无记账记录</p>
                <p style="font-size: 0.9rem; margin-top: 10px;">点击上方表单添加第一条记录</p>
            </div>
        `;
        return;
    }
    
    container.innerHTML = records.map(record => {
        const typeInfo = TYPE_MAP[record.type];
        const amount = record.type === 'expense' ? `-${formatMoney(record.amount)}` : 
                      record.type === 'income' ? `+${formatMoney(record.amount)}` : 
                      formatMoney(record.amount);
        
        return `
            <div class="record-item ${record.type}">
                <div class="record-icon">${typeInfo.icon}</div>
                <div class="record-info">
                    <div class="record-category">
                        ${record.category_name}${record.subcategory_name ? ' > ' + record.subcategory_name : ''}
                    </div>
                    ${record.note ? `<div class="record-note">${record.note}</div>` : ''}
                    <div class="record-date">${formatDate(record.date)}</div>
                </div>
                <div class="record-amount">${amount}</div>
                <div class="record-actions">
                    <button class="btn btn-edit" onclick="editRecord(${record.id})">编辑</button>
                    <button class="btn btn-danger" onclick="deleteRecord(${record.id})">删除</button>
                </div>
            </div>
        `;
    }).join('');
}

// 添加记录
async function handleAddRecord(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const data = {
        type: formData.get('type'),
        category_id: parseInt(formData.get('category_id')),
        subcategory_id: formData.get('subcategory_id') ? parseInt(formData.get('subcategory_id')) : null,
        amount: parseFloat(formData.get('amount')),
        date: formData.get('date'),
        note: formData.get('note')
    };
    
    try {
        const response = await fetch(`${API_BASE}/records`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (result.Errno === 0) {
            showSuccess('添加成功');
            document.getElementById('recordForm').reset();
            initDateInputs();
            loadRecords();
            loadStatistics();
        } else {
            showError('添加失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('添加失败: ' + error.message);
    }
}

// 编辑记录
async function editRecord(id) {
    const record = records.find(r => r.id === id);
    if (!record) return;
    
    document.getElementById('editId').value = record.id;
    document.getElementById('editType').value = record.type;
    document.getElementById('editAmount').value = record.amount;
    document.getElementById('editDate').value = record.date;
    document.getElementById('editNote').value = record.note || '';
    
    // 更新分类选项
    updateEditCategoryOptions();
    document.getElementById('editCategory').value = record.category_id;
    
    // 更新子分类选项
    updateEditSubcategoryOptions();
    if (record.subcategory_id) {
        document.getElementById('editSubcategory').value = record.subcategory_id;
    }
    
    document.getElementById('editModal').classList.add('active');
}

// 更新记录
async function handleUpdateRecord(e) {
    e.preventDefault();
    
    const id = document.getElementById('editId').value;
    const formData = new FormData(e.target);
    const data = {
        type: formData.get('type'),
        category_id: parseInt(formData.get('category_id')),
        subcategory_id: formData.get('subcategory_id') ? parseInt(formData.get('subcategory_id')) : null,
        amount: parseFloat(formData.get('amount')),
        date: formData.get('date'),
        note: formData.get('note')
    };
    
    try {
        const response = await fetch(`${API_BASE}/records/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (result.Errno === 0) {
            showSuccess('更新成功');
            closeEditModal();
            loadRecords();
            loadStatistics();
        } else {
            showError('更新失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('更新失败: ' + error.message);
    }
}

// 删除记录
async function deleteRecord(id) {
    if (!confirm('确定要删除这条记录吗？')) return;
    
    try {
        const response = await fetch(`${API_BASE}/records/${id}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        
        if (result.Errno === 0) {
            showSuccess('删除成功');
            loadRecords();
            loadStatistics();
        } else {
            showError('删除失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('删除失败: ' + error.message);
    }
}

// 关闭编辑模态框
function closeEditModal() {
    document.getElementById('editModal').classList.remove('active');
}

// 筛选
function handleFilter() {
    const filters = {
        type: document.getElementById('filterType').value,
        categoryId: document.getElementById('filterCategory').value,
        startDate: document.getElementById('filterStartDate').value,
        endDate: document.getElementById('filterEndDate').value,
        keyword: document.getElementById('filterKeyword').value
    };
    
    loadRecords(filters);
}

// 清空筛选
function clearFilter() {
    document.getElementById('filterType').value = '';
    document.getElementById('filterCategory').value = '';
    document.getElementById('filterStartDate').value = '';
    document.getElementById('filterEndDate').value = '';
    document.getElementById('filterKeyword').value = '';
    loadRecords();
}

// ==================== 分类管理功能 ====================

// 渲染分类管理页面
function renderCategoryManagement() {
    const types = ['expense', 'income', 'transfer', 'loan'];
    
    types.forEach(type => {
        const container = document.getElementById(type + 'Categories');
        const typeCategories = categories.filter(cat => cat.type === type);
        
        if (typeCategories.length === 0) {
            container.innerHTML = '<p style="color: #999; padding: 10px;">暂无分类</p>';
            return;
        }
        
        container.innerHTML = typeCategories.map(cat => `
            <div class="category-item">
                <div class="category-header">
                    <span class="category-name">${cat.name}</span>
                    <div class="category-actions">
                        <button class="btn btn-edit btn-small" onclick="openEditCategoryModal(${cat.id}, '${cat.name}', '${cat.type}')">编辑</button>
                        <button class="btn btn-danger btn-small" onclick="deleteCategory(${cat.id})">删除</button>
                    </div>
                </div>
                <div class="subcategory-list">
                    ${cat.subcategories && cat.subcategories.length > 0 
                        ? cat.subcategories.map(sub => `
                            <div class="subcategory-item">
                                <span class="subcategory-name">${sub.name}</span>
                                <div class="subcategory-actions">
                                    <button class="btn-icon edit" onclick="openEditSubcategoryModal(${sub.id}, '${sub.name}')">✏️</button>
                                    <button class="btn-icon delete" onclick="deleteSubcategory(${sub.id})">🗑️</button>
                                </div>
                            </div>
                        `).join('')
                        : ''
                    }
                    <button class="add-subcategory-btn" onclick="openAddSubcategoryModal(${cat.id})">+ 添加子分类</button>
                </div>
            </div>
        `).join('');
    });
}

// 添加一级分类
async function addCategory(type) {
    const inputId = 'new' + type.charAt(0).toUpperCase() + type.slice(1) + 'Category';
    const input = document.getElementById(inputId);
    const name = input.value.trim();
    
    if (!name) {
        showError('请输入分类名称');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/categories`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, type })
        });
        
        const result = await response.json();
        
        if (result.Errno === 0) {
            showSuccess('添加成功');
            input.value = '';
            await loadCategories();
            renderCategoryManagement();
        } else {
            showError('添加失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('添加失败: ' + error.message);
    }
}

// 打开编辑分类模态框
function openEditCategoryModal(id, name, type) {
    document.getElementById('editCategoryId').value = id;
    document.getElementById('editCategoryType').value = type;
    document.getElementById('editCategoryName').value = name;
    document.getElementById('editCategoryModal').classList.add('active');
}

// 关闭编辑分类模态框
function closeEditCategoryModal() {
    document.getElementById('editCategoryModal').classList.remove('active');
}

// 更新分类
async function handleUpdateCategory(e) {
    e.preventDefault();
    
    const id = document.getElementById('editCategoryId').value;
    const name = document.getElementById('editCategoryName').value.trim();
    const type = document.getElementById('editCategoryType').value;
    
    if (!name) {
        showError('请输入分类名称');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/categories/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, type })
        });
        
        const result = await response.json();
        
        if (result.Errno === 0) {
            showSuccess('更新成功');
            closeEditCategoryModal();
            await loadCategories();
            renderCategoryManagement();
            updateFilterCategoryOptions();
        } else {
            showError('更新失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('更新失败: ' + error.message);
    }
}

// 删除分类
async function deleteCategory(id) {
    if (!confirm('确定要删除这个分类吗？该分类下的所有子分类也会被删除。')) return;
    
    try {
        const response = await fetch(`${API_BASE}/categories/${id}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        
        if (result.Errno === 0) {
            showSuccess('删除成功');
            await loadCategories();
            renderCategoryManagement();
            updateFilterCategoryOptions();
        } else {
            showError('删除失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('删除失败: ' + error.message);
    }
}

// 打开添加子分类模态框
function openAddSubcategoryModal(categoryId) {
    document.getElementById('addSubcategoryParentId').value = categoryId;
    document.getElementById('addSubcategoryName').value = '';
    document.getElementById('addSubcategoryModal').classList.add('active');
}

// 关闭添加子分类模态框
function closeAddSubcategoryModal() {
    document.getElementById('addSubcategoryModal').classList.remove('active');
}

// 添加子分类
async function handleAddSubcategory(e) {
    e.preventDefault();
    
    const categoryId = document.getElementById('addSubcategoryParentId').value;
    const name = document.getElementById('addSubcategoryName').value.trim();
    
    if (!name) {
        showError('请输入子分类名称');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/categories/subcategories`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ category_id: parseInt(categoryId), name })
        });
        
        const result = await response.json();
        
        if (result.Errno === 0) {
            showSuccess('添加成功');
            closeAddSubcategoryModal();
            await loadCategories();
            renderCategoryManagement();
        } else {
            showError('添加失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('添加失败: ' + error.message);
    }
}

// 打开编辑子分类模态框
function openEditSubcategoryModal(id, name) {
    document.getElementById('editSubcategoryId').value = id;
    document.getElementById('editSubcategoryName').value = name;
    document.getElementById('editSubcategoryModal').classList.add('active');
}

// 关闭编辑子分类模态框
function closeEditSubcategoryModal() {
    document.getElementById('editSubcategoryModal').classList.remove('active');
}

// 更新子分类
async function handleUpdateSubcategory(e) {
    e.preventDefault();
    
    const id = document.getElementById('editSubcategoryId').value;
    const name = document.getElementById('editSubcategoryName').value.trim();
    
    if (!name) {
        showError('请输入子分类名称');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/categories/subcategories/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name })
        });
        
        const result = await response.json();
        
        if (result.Errno === 0) {
            showSuccess('更新成功');
            closeEditSubcategoryModal();
            await loadCategories();
            renderCategoryManagement();
        } else {
            showError('更新失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('更新失败: ' + error.message);
    }
}

// 删除子分类
async function deleteSubcategory(id) {
    if (!confirm('确定要删除这个子分类吗？')) return;
    
    try {
        const response = await fetch(`${API_BASE}/categories/subcategories/${id}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        
        if (result.Errno === 0) {
            showSuccess('删除成功');
            await loadCategories();
            renderCategoryManagement();
        } else {
            showError('删除失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('删除失败: ' + error.message);
    }
}

// ==================== 工具函数 ====================

// 格式化金额
function formatMoney(amount) {
    return '¥' + parseFloat(amount).toFixed(2);
}

// 格式化日期
function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
    });
}

// 显示成功消息
function showSuccess(message) {
    // 可以替换为更友好的提示方式
    alert(message);
}

// 显示错误消息
function showError(message) {
    // 可以替换为更友好的提示方式
    alert(message);
}
