<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	
	const API_URL = 'http://localhost:3000';
	
	interface Category {
		id: number;
		profile_id: number;
		name: string;
		icon: string;
		color: string;
		position: number;
		is_expanded: boolean;
	}
	
	let categories = $state<Category[]>([]);
	let loading = $state(true);
	let profileId = $state<number | null>(null);
	
	// Modal state
	let showModal = $state(false);
	let editingCategory = $state<Partial<Category> | null>(null);
	let isNew = $state(false);
	let saving = $state(false);
	
	onMount(async () => {
		const id = $page.url.searchParams.get('profile');
		if (id) {
			profileId = parseInt(id);
			await loadCategories();
		}
		loading = false;
	});
	
	async function loadCategories() {
		if (!profileId) return;
		const token = localStorage.getItem('access_token');
		const res = await fetch(`${API_URL}/api/v1/profiles/${profileId}/categories`, {
			headers: { Authorization: `Bearer ${token}` }
		});
		if (res.ok) {
			const data = await res.json();
			categories = data.data || [];
		}
	}
	
	function openNewCategory() {
		editingCategory = {
			name: '',
			icon: 'bi-folder',
			color: '#667eea',
			is_expanded: true
		};
		isNew = true;
		showModal = true;
	}
	
	function openEditCategory(cat: Category) {
		editingCategory = { ...cat };
		isNew = false;
		showModal = true;
	}
	
	function closeModal() {
		showModal = false;
		editingCategory = null;
	}
	
	async function saveCategory() {
		if (!editingCategory || !profileId) return;
		saving = true;
		
		const token = localStorage.getItem('access_token');
		
		if (isNew) {
			const res = await fetch(`${API_URL}/api/v1/profiles/${profileId}/categories`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${token}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(editingCategory)
			});
			if (res.ok) {
				const data = await res.json();
				if (data.data) {
					categories = [...categories, data.data];
				}
			}
		} else if (editingCategory.id) {
			await fetch(`${API_URL}/api/v1/categories/${editingCategory.id}`, {
				method: 'PUT',
				headers: {
					Authorization: `Bearer ${token}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(editingCategory)
			});
			categories = categories.map(c => c.id === editingCategory!.id ? { ...c, ...editingCategory } as Category : c);
		}
		
		saving = false;
		closeModal();
	}
	
	async function deleteCategory(id: number) {
		if (!confirm('Delete this category? Links in it will be moved to uncategorized.')) return;
		
		const token = localStorage.getItem('access_token');
		await fetch(`${API_URL}/api/v1/categories/${id}`, {
			method: 'DELETE',
			headers: { Authorization: `Bearer ${token}` }
		});
		categories = categories.filter(c => c.id !== id);
	}
	
	const iconOptions = [
		'bi-folder', 'bi-collection', 'bi-bookmark', 'bi-star', 'bi-heart',
		'bi-lightning', 'bi-globe', 'bi-briefcase', 'bi-music-note', 'bi-camera',
		'bi-cart', 'bi-chat', 'bi-book', 'bi-calendar', 'bi-gear'
	];
	
	const colorOptions = [
		'#667eea', '#764ba2', '#f093fb', '#f5576c', '#4facfe',
		'#00f2fe', '#43e97b', '#38f9d7', '#ffecd2', '#fcb69f',
		'#ff6b6b', '#feca57', '#48dbfb', '#ff9ff3', '#54a0ff'
	];
</script>

<svelte:head>
	<title>Categories - LinkMy</title>
</svelte:head>

<div class="categories-page">
	<header class="page-header">
		<div class="header-left">
			<h1>Categories</h1>
			<p>Organize your links into folders</p>
		</div>
		<div class="header-actions">
			<button class="btn btn-primary" onclick={openNewCategory}>
				<i class="bi bi-plus-lg"></i>
				New Category
			</button>
		</div>
	</header>
	
	{#if loading}
		<div class="loading">Loading...</div>
	{:else if !profileId}
		<div class="card empty-state">
			<i class="bi bi-exclamation-triangle"></i>
			<h3>No profile selected</h3>
			<p>Please select a profile first</p>
			<a href="/dashboard/profiles" class="btn btn-primary">Go to Profiles</a>
		</div>
	{:else if categories.length === 0}
		<div class="card empty-state">
			<i class="bi bi-folder"></i>
			<h3>No categories yet</h3>
			<p>Create categories to organize your links</p>
			<button class="btn btn-primary" onclick={openNewCategory}>
				<i class="bi bi-plus-lg"></i>
				Create Category
			</button>
		</div>
	{:else}
		<div class="categories-list">
			{#each categories as cat, i (cat.id)}
				<div class="category-card card">
					<div class="category-icon" style="background-color: {cat.color}">
						<i class="bi {cat.icon}"></i>
					</div>
					<div class="category-info">
						<h3>{cat.name}</h3>
						<span class="position">Position: {cat.position + 1}</span>
					</div>
					<div class="category-status">
						{#if cat.is_expanded}
							<span class="badge expanded">Expanded</span>
						{:else}
							<span class="badge collapsed">Collapsed</span>
						{/if}
					</div>
					<div class="category-actions">
						<button class="btn-icon" onclick={() => openEditCategory(cat)} title="Edit">
							<i class="bi bi-pencil"></i>
						</button>
						<button class="btn-icon danger" onclick={() => deleteCategory(cat.id)} title="Delete">
							<i class="bi bi-trash"></i>
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Category Modal -->
{#if showModal && editingCategory}
	<div class="modal-overlay" onclick={closeModal}>
		<div class="modal card" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2>{isNew ? 'New Category' : 'Edit Category'}</h2>
				<button class="btn-icon" onclick={closeModal}>
					<i class="bi bi-x-lg"></i>
				</button>
			</div>
			
			<form onsubmit={(e) => { e.preventDefault(); saveCategory(); }}>
				<div class="form-group">
					<label for="name">Name</label>
					<input
						type="text"
						id="name"
						bind:value={editingCategory.name}
						placeholder="Category name"
						required
					/>
				</div>
				
				<div class="form-group">
					<label>Icon</label>
					<div class="icon-grid">
						{#each iconOptions as icon}
							<button 
								type="button"
								class="icon-btn"
								class:selected={editingCategory.icon === icon}
								onclick={() => editingCategory!.icon = icon}
							>
								<i class="bi {icon}"></i>
							</button>
						{/each}
					</div>
				</div>
				
				<div class="form-group">
					<label>Color</label>
					<div class="color-grid">
						{#each colorOptions as color}
							<button 
								type="button"
								class="color-btn"
								class:selected={editingCategory.color === color}
								style="background-color: {color}"
								onclick={() => editingCategory!.color = color}
							></button>
						{/each}
					</div>
				</div>
				
				<div class="form-group">
					<label class="checkbox-label">
						<input type="checkbox" bind:checked={editingCategory.is_expanded} />
						<span>Expanded by default</span>
					</label>
				</div>
				
				<div class="form-actions">
					<button type="button" class="btn btn-secondary" onclick={closeModal}>Cancel</button>
					<button type="submit" class="btn btn-primary" disabled={saving}>
						{saving ? 'Saving...' : 'Save'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<style>
	.categories-page {
		max-width: 800px;
	}
	
	.page-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-2xl);
	}
	
	.page-header h1 {
		font-size: 2rem;
		margin-bottom: var(--space-xs);
	}
	
	.page-header p {
		color: var(--color-text-secondary);
	}
	
	.loading {
		text-align: center;
		padding: var(--space-3xl);
		color: var(--color-text-muted);
	}
	
	.empty-state {
		text-align: center;
		padding: var(--space-3xl);
	}
	
	.empty-state i {
		font-size: 3rem;
		color: var(--color-text-muted);
		margin-bottom: var(--space-lg);
	}
	
	.empty-state h3 {
		margin-bottom: var(--space-sm);
	}
	
	.empty-state p {
		color: var(--color-text-secondary);
		margin-bottom: var(--space-xl);
	}
	
	.categories-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
	}
	
	.category-card {
		display: flex;
		align-items: center;
		gap: var(--space-lg);
		padding: var(--space-lg);
	}
	
	.category-icon {
		width: 48px;
		height: 48px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: var(--radius-lg);
		color: white;
		font-size: 1.25rem;
	}
	
	.category-info {
		flex: 1;
	}
	
	.category-info h3 {
		font-size: 1.125rem;
		margin-bottom: var(--space-xs);
	}
	
	.position {
		font-size: 0.75rem;
		color: var(--color-text-muted);
	}
	
	.badge {
		padding: 2px 8px;
		border-radius: 999px;
		font-size: 0.675rem;
		font-weight: 500;
		text-transform: uppercase;
	}
	
	.badge.expanded { background: rgba(16, 185, 129, 0.2); color: #10b981; }
	.badge.collapsed { background: rgba(107, 114, 128, 0.2); color: #6b7280; }
	
	.category-actions {
		display: flex;
		gap: var(--space-xs);
	}
	
	.btn-icon {
		width: 36px;
		height: 36px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		transition: all var(--transition-fast);
	}
	
	.btn-icon:hover {
		background: var(--color-bg-tertiary);
		color: var(--color-text);
	}
	
	.btn-icon.danger:hover {
		background: rgba(239, 68, 68, 0.1);
		color: #ef4444;
	}
	
	/* Modal */
	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-lg);
		z-index: 100;
	}
	
	.modal {
		width: 100%;
		max-width: 480px;
		max-height: 90vh;
		overflow-y: auto;
	}
	
	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-xl);
	}
	
	.modal-header h2 {
		font-size: 1.25rem;
	}
	
	.form-group {
		margin-bottom: var(--space-lg);
	}
	
	.form-group label {
		display: block;
		font-size: 0.875rem;
		font-weight: 500;
		margin-bottom: var(--space-sm);
		color: var(--color-text-secondary);
	}
	
	.form-group input[type="text"] {
		width: 100%;
		padding: var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text);
		font-size: 1rem;
	}
	
	.icon-grid {
		display: grid;
		grid-template-columns: repeat(5, 1fr);
		gap: var(--space-xs);
	}
	
	.icon-btn {
		width: 48px;
		height: 48px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-bg-tertiary);
		border: 2px solid transparent;
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		font-size: 1.25rem;
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.icon-btn:hover {
		border-color: var(--color-border);
	}
	
	.icon-btn.selected {
		border-color: var(--color-primary);
		color: var(--color-primary);
	}
	
	.color-grid {
		display: grid;
		grid-template-columns: repeat(5, 1fr);
		gap: var(--space-xs);
	}
	
	.color-btn {
		width: 48px;
		height: 48px;
		border: 3px solid transparent;
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.color-btn:hover {
		transform: scale(1.1);
	}
	
	.color-btn.selected {
		border-color: white;
		box-shadow: 0 0 0 2px var(--color-primary);
	}
	
	.checkbox-label {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		cursor: pointer;
	}
	
	.checkbox-label input {
		width: 18px;
		height: 18px;
		accent-color: var(--color-primary);
	}
	
	.form-actions {
		display: flex;
		gap: var(--space-md);
		justify-content: flex-end;
		margin-top: var(--space-xl);
	}
</style>
