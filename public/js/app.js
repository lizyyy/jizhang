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
});

// 初始化日期输入框
function initDateInputs() {
    const today = new Date().toISOString().split('T')[0];
    document.getElementById('date').value = today;
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
            e.target.reset();
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
    
    // 更新分类选项
    updateEditCategoryOptions();
    
    // 设置分类值
    document.getElementById('editCategory').value = record.category_id;
    
    // 更新子分类选项
    updateEditSubcategoryOptions();
    
    // 设置其他值
    setTimeout(() => {
        document.getElementById('editSubcategory').value = record.subcategory_id || '';
    }, 0);
    
    document.getElementById('editAmount').value = record.amount;
    document.getElementById('editDate').value = record.date;
    document.getElementById('editNote').value = record.note || '';
    
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

// 格式化金额
function formatMoney(amount) {
    return '¥' + parseFloat(amount).toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}

// 格式化日期
function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

// 显示成功消息
function showSuccess(message) {
    // 简单的 alert，可以替换为更优雅的提示
    alert(message);
}

// 显示错误消息
function showError(message) {
    alert('错误: ' + message);
}

// ==================== 分类管理功能 ====================

let currentCategoryType = 'expense';

// 绑定分类管理事件
function bindCategoryEvents() {
    // 添加分类按钮
    document.getElementById('addCategoryBtn').addEventListener('click', () => {
        openCategoryModal('create');
    });

    // 分类类型标签切换
    document.querySelectorAll('.category-tab').forEach(tab => {
        tab.addEventListener('click', (e) => {
            document.querySelectorAll('.category-tab').forEach(t => t.classList.remove('active'));
            e.target.classList.add('active');
            currentCategoryType = e.target.dataset.type;
            renderCategoryList();
        });
    });

    // 分类表单提交
    document.getElementById('categoryForm').addEventListener('submit', handleCategorySubmit);

    // 关闭分类模态框
    document.getElementById('closeCategoryModal').addEventListener('click', closeCategoryModal);
    document.getElementById('cancelCategory').addEventListener('click', closeCategoryModal);

    // 点击模态框外部关闭
    document.getElementById('categoryModal').addEventListener('click', (e) => {
        if (e.target.id === 'categoryModal') {
            closeCategoryModal();
        }
    });
}

// 渲染分类列表
function renderCategoryList() {
    const container = document.getElementById('categoryList');
    const filteredCategories = categories.filter(cat => cat.type === currentCategoryType);

    if (filteredCategories.length === 0) {
        container.innerHTML = `
            <div class="category-empty">
                <div class="category-empty-icon">📂</div>
                <p>暂无${TYPE_MAP[currentCategoryType].label}分类</p>
                <p style="font-size: 0.9rem; margin-top: 10px;">点击上方按钮添加新分类</p>
            </div>
        `;
        return;
    }

    container.innerHTML = filteredCategories.map(category => {
        const typeInfo = TYPE_MAP[category.type];
        return `
            <div class="category-item ${category.type}" data-id="${category.id}">
                <div class="category-header">
                    <div class="category-info">
                        <div class="category-icon">${typeInfo.icon}</div>
                        <div class="category-name">${category.name}</div>
                    </div>
                    <div class="category-actions">
                        <button class="btn btn-edit" onclick="editCategory(${category.id})">编辑</button>
                        <button class="btn btn-danger" onclick="deleteCategory(${category.id})">删除</button>
                    </div>
                </div>
                ${renderSubcategories(category)}
                <button class="add-subcategory-btn" onclick="openSubcategoryModal(${category.id})">+ 添加子分类</button>
            </div>
        `;
    }).join('');
}

// 渲染子分类列表
function renderSubcategories(category) {
    if (!category.subcategories || category.subcategories.length === 0) {
        return '';
    }

    return `
        <div class="subcategory-list">
            ${category.subcategories.map(sub => `
                <div class="subcategory-item" data-id="${sub.id}">
                    <div class="subcategory-name">${sub.name}</div>
                    <div class="subcategory-actions">
                        <button class="btn btn-edit" onclick="editSubcategory(${category.id}, ${sub.id}, '${sub.name}')">编辑</button>
                        <button class="btn btn-danger" onclick="deleteSubcategory(${sub.id})">删除</button>
                    </div>
                </div>
            `).join('')}
        </div>
    `;
}

// 打开分类模态框
function openCategoryModal(action, categoryId = null) {
    const modal = document.getElementById('categoryModal');
    const title = document.getElementById('categoryModalTitle');
    const form = document.getElementById('categoryForm');
    const typeGroup = document.getElementById('categoryTypeGroup');

    form.reset();
    document.getElementById('categoryAction').value = action;
    document.getElementById('categoryId').value = categoryId || '';
    document.getElementById('parentCategoryId').value = '';

    if (action === 'create') {
        title.textContent = '添加分类';
        typeGroup.style.display = 'block';
        document.getElementById('categoryType').value = currentCategoryType;
    } else if (action === 'edit') {
        title.textContent = '编辑分类';
        typeGroup.style.display = 'block';
        // 填充表单数据
        const category = categories.find(c => c.id === categoryId);
        if (category) {
            document.getElementById('categoryName').value = category.name;
            document.getElementById('categoryType').value = category.type;
        }
    } else if (action === 'createSub') {
        title.textContent = '添加子分类';
        typeGroup.style.display = 'none';
    }

    modal.classList.add('active');
}

// 打开子分类模态框
function openSubcategoryModal(categoryId) {
    const modal = document.getElementById('categoryModal');
    const title = document.getElementById('categoryModalTitle');
    const form = document.getElementById('categoryForm');
    const typeGroup = document.getElementById('categoryTypeGroup');

    form.reset();
    document.getElementById('categoryAction').value = 'createSub';
    document.getElementById('parentCategoryId').value = categoryId;
    document.getElementById('categoryId').value = '';

    title.textContent = '添加子分类';
    typeGroup.style.display = 'none';

    modal.classList.add('active');
}

// 编辑分类
function editCategory(categoryId) {
    openCategoryModal('edit', categoryId);
}

// 编辑子分类
function editSubcategory(categoryId, subId, name) {
    const modal = document.getElementById('categoryModal');
    const title = document.getElementById('categoryModalTitle');
    const form = document.getElementById('categoryForm');
    const typeGroup = document.getElementById('categoryTypeGroup');

    form.reset();
    document.getElementById('categoryAction').value = 'editSub';
    document.getElementById('categoryId').value = subId;
    document.getElementById('parentCategoryId').value = categoryId;
    document.getElementById('categoryName').value = name;

    title.textContent = '编辑子分类';
    typeGroup.style.display = 'none';

    modal.classList.add('active');
}

// 关闭分类模态框
function closeCategoryModal() {
    document.getElementById('categoryModal').classList.remove('active');
}

// 处理分类表单提交
async function handleCategorySubmit(e) {
    e.preventDefault();

    const action = document.getElementById('categoryAction').value;
    const name = document.getElementById('categoryName').value;
    const type = document.getElementById('categoryType').value;
    const categoryId = document.getElementById('categoryId').value;
    const parentId = document.getElementById('parentCategoryId').value;

    let url = `${API_BASE}/categories`;
    let method = 'POST';
    let body = {};

    switch (action) {
        case 'create':
            body = { name, type };
            break;
        case 'edit':
            url += `/${categoryId}`;
            method = 'PUT';
            body = { name, type };
            break;
        case 'createSub':
            url += '/subcategories';
            body = { category_id: parseInt(parentId), name };
            break;
        case 'editSub':
            url += `/subcategories/${categoryId}`;
            method = 'PUT';
            body = { name };
            break;
    }

    try {
        const response = await fetch(url, {
            method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });

        const result = await response.json();

        if (result.Errno === 0) {
            showSuccess('保存成功');
            closeCategoryModal();
            await loadCategories();
            renderCategoryList();
            updateCategoryOptions();
            updateFilterCategoryOptions();
        } else {
            showError('保存失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('保存失败: ' + error.message);
    }
}

// 删除分类
async function deleteCategory(categoryId) {
    if (!confirm('确定要删除这个分类吗？该分类下的所有子分类也将被删除。')) return;

    try {
        const response = await fetch(`${API_BASE}/categories/${categoryId}`, {
            method: 'DELETE'
        });

        const result = await response.json();

        if (result.Errno === 0) {
            showSuccess('删除成功');
            await loadCategories();
            renderCategoryList();
            updateCategoryOptions();
            updateFilterCategoryOptions();
        } else {
            showError('删除失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('删除失败: ' + error.message);
    }
}

// 删除子分类
async function deleteSubcategory(subId) {
    if (!confirm('确定要删除这个子分类吗？')) return;

    try {
        const response = await fetch(`${API_BASE}/categories/subcategories/${subId}`, {
            method: 'DELETE'
        });

        const result = await response.json();

        if (result.Errno === 0) {
            showSuccess('删除成功');
            await loadCategories();
            renderCategoryList();
            updateCategoryOptions();
        } else {
            showError('删除失败: ' + result.Errmsg);
        }
    } catch (error) {
        showError('删除失败: ' + error.message);
    }
}

// 修改loadCategories函数，在加载完成后渲染分类列表
const originalLoadCategories = loadCategories;
loadCategories = async function() {
    await originalLoadCategories();
    if (typeof renderCategoryList === 'function') {
        renderCategoryList();
    }
}

// 修改bindEvents函数，添加分类管理事件绑定
const originalBindEvents = bindEvents;
bindEvents = function() {
    originalBindEvents();
    if (typeof bindCategoryEvents === 'function') {
        bindCategoryEvents();
    }
}
