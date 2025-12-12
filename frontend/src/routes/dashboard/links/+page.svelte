<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { links as linksApi, profiles, type Link, type Profile } from '$lib/api';
	import IconPicker from '$lib/components/IconPicker.svelte';
	
	let userProfiles = $state<Profile[]>([]);
	let selectedProfileId = $state<number | null>(null);
	let userLinks = $state<Link[]>([]);
	let loading = $state(true);
	let saving = $state(false);
	
	// Edit modal state
	let showModal = $state(false);
	let editingLink = $state<Partial<Link> | null>(null);
	let isNewLink = $state(false);
	
	onMount(async () => {
		// Load profiles
		const res = await profiles.getAll();
		if (res.data && res.data.length > 0) {
			userProfiles = res.data;
			// Check URL for profile param
			const urlProfile = parseInt($page.url.searchParams.get('profile') || '');
			selectedProfileId = urlProfile || res.data[0].id;
			await loadLinks();
		}
		loading = false;
	});
	
	async function loadLinks() {
		if (!selectedProfileId) return;
		const res = await linksApi.getByProfile(selectedProfileId);
		if (res.data) {
			userLinks = res.data;
		}
	}
	
	async function handleProfileChange(e: Event) {
		const select = e.target as HTMLSelectElement;
		selectedProfileId = parseInt(select.value);
		await loadLinks();
	}
	
	function openNewLink() {
		editingLink = {
			title: '',
			url: '',
			icon: 'bi-link-45deg',
			is_active: true
		};
		isNewLink = true;
		showModal = true;
	}
	
	function openEditLink(link: Link) {
		editingLink = { ...link };
		isNewLink = false;
		showModal = true;
	}
	
	function closeModal() {
		showModal = false;
		editingLink = null;
	}
	
	async function saveLink() {
		if (!editingLink || !selectedProfileId) return;
		saving = true;
		
		if (isNewLink) {
			const res = await linksApi.create(selectedProfileId, editingLink);
			if (res.data) {
				userLinks = [...userLinks, res.data];
			}
		} else if (editingLink.id) {
			const res = await linksApi.update(editingLink.id, editingLink);
			if (res.data) {
				userLinks = userLinks.map(l => l.id === res.data!.id ? res.data! : l);
			}
		}
		
		saving = false;
		closeModal();
	}
	
	async function deleteLink(id: number) {
		if (!confirm('Delete this link?')) return;
		await linksApi.delete(id);
		userLinks = userLinks.filter(l => l.id !== id);
	}
	
	async function toggleActive(link: Link) {
		await linksApi.update(link.id, { is_active: !link.is_active });
		userLinks = userLinks.map(l => l.id === link.id ? { ...l, is_active: !l.is_active } : l);
	}
</script>

<svelte:head>
	<title>Links - LinkMy</title>
</svelte:head>

<div class="links-page">
	<header class="page-header">
		<div class="header-left">
			<h1>Links</h1>
			<p>Manage your links</p>
		</div>
		<div class="header-actions">
			{#if userProfiles.length > 1}
				<select class="profile-select" onchange={handleProfileChange} value={selectedProfileId}>
					{#each userProfiles as profile}
						<option value={profile.id}>{profile.name}</option>
					{/each}
				</select>
			{/if}
			<button class="btn btn-primary" onclick={openNewLink}>
				<i class="bi bi-plus-lg"></i>
				Add Link
			</button>
		</div>
	</header>
	
	{#if loading}
		<div class="loading-state">
			<div class="loader"></div>
		</div>
	{:else if userLinks.length === 0}
		<div class="empty-state card">
			<i class="bi bi-link-45deg"></i>
			<h3>No links yet</h3>
			<p>Add your first link to get started</p>
			<button class="btn btn-primary" onclick={openNewLink}>
				<i class="bi bi-plus-lg"></i>
				Add Link
			</button>
		</div>
	{:else}
		<div class="links-list">
			{#each userLinks as link, i (link.id)}
				<div class="link-item card" class:inactive={!link.is_active}>
					<div class="link-drag">
						<i class="bi bi-grip-vertical"></i>
					</div>
					<div class="link-icon">
						<i class="bi {link.icon}"></i>
					</div>
					<div class="link-info">
						<h3>{link.title}</h3>
						<a href={link.url} target="_blank" class="link-url">{link.url}</a>
					</div>
					<div class="link-stats">
						<span class="clicks">
							<i class="bi bi-cursor"></i>
							{link.clicks}
						</span>
					</div>
					<div class="link-actions">
						<button 
							class="btn-icon" 
							onclick={() => toggleActive(link)}
							title={link.is_active ? 'Disable' : 'Enable'}
						>
							<i class="bi {link.is_active ? 'bi-eye' : 'bi-eye-slash'}"></i>
						</button>
						<button class="btn-icon" onclick={() => openEditLink(link)} title="Edit">
							<i class="bi bi-pencil"></i>
						</button>
						<button class="btn-icon danger" onclick={() => deleteLink(link.id)} title="Delete">
							<i class="bi bi-trash"></i>
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Edit/Create Modal -->
{#if showModal && editingLink}
	<div class="modal-overlay" onclick={closeModal}>
		<div class="modal card" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2>{isNewLink ? 'Add Link' : 'Edit Link'}</h2>
				<button class="btn-icon" onclick={closeModal}>
					<i class="bi bi-x-lg"></i>
				</button>
			</div>
			
			<form onsubmit={(e) => { e.preventDefault(); saveLink(); }}>
				<div class="form-group">
					<label for="title">Title</label>
					<input
						type="text"
						id="title"
						bind:value={editingLink.title}
						placeholder="My Website"
						required
					/>
				</div>
				
				<div class="form-group">
					<label for="url">URL</label>
					<input
						type="url"
						id="url"
						bind:value={editingLink.url}
						placeholder="https://example.com"
						required
					/>
				</div>
				
				<div class="form-group">
					<label>Icon</label>
					<IconPicker bind:value={editingLink.icon} />
				
				<div class="form-actions">
					<button type="button" class="btn btn-secondary" onclick={closeModal}>
						Cancel
					</button>
					<button type="submit" class="btn btn-primary" disabled={saving}>
						{#if saving}
							<i class="bi bi-arrow-repeat spin"></i>
							Saving...
						{:else}
							<i class="bi bi-check-lg"></i>
							Save
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<style>
	.links-page {
		max-width: 900px;
	}
	
	.page-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-2xl);
		flex-wrap: wrap;
		gap: var(--space-lg);
	}
	
	.page-header h1 {
		font-size: 2rem;
		margin-bottom: var(--space-xs);
	}
	
	.page-header p {
		color: var(--color-text-secondary);
	}
	
	.header-actions {
		display: flex;
		gap: var(--space-md);
	}
	
	.profile-select {
		padding: var(--space-md) var(--space-lg);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text);
		font-size: 0.875rem;
	}
	
	/* Links List */
	.links-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
	}
	
	.link-item {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		padding: var(--space-lg);
		transition: all var(--transition-fast);
	}
	
	.link-item.inactive {
		opacity: 0.5;
	}
	
	.link-drag {
		color: var(--color-text-muted);
		cursor: grab;
	}
	
	.link-icon {
		width: 44px;
		height: 44px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--gradient-primary);
		border-radius: var(--radius-md);
		color: white;
		font-size: 1.25rem;
	}
	
	.link-info {
		flex: 1;
		min-width: 0;
	}
	
	.link-info h3 {
		font-size: 1rem;
		font-weight: 500;
		margin-bottom: var(--space-xs);
	}
	
	.link-url {
		font-size: 0.75rem;
		color: var(--color-text-muted);
		text-overflow: ellipsis;
		overflow: hidden;
		white-space: nowrap;
		display: block;
	}
	
	.link-stats {
		display: flex;
		gap: var(--space-lg);
	}
	
	.clicks {
		display: flex;
		align-items: center;
		gap: var(--space-xs);
		font-size: 0.875rem;
		color: var(--color-text-secondary);
	}
	
	.link-actions {
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
	
	/* Empty State */
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
	
	/* Loading */
	.loading-state {
		display: flex;
		justify-content: center;
		padding: var(--space-3xl);
	}
	
	.loader {
		width: 40px;
		height: 40px;
		border: 3px solid var(--color-border);
		border-top-color: var(--color-primary);
		border-radius: 50%;
		animation: spin 1s linear infinite;
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
	
	.form-group input,
	.form-group select {
		width: 100%;
		padding: var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text);
		font-size: 1rem;
	}
	
	.form-group input:focus,
	.form-group select:focus {
		outline: none;
		border-color: var(--color-primary);
	}
	
	.icon-select {
		display: flex;
		gap: var(--space-md);
	}
	
	.icon-select select {
		flex: 1;
	}
	
	.icon-preview {
		width: 48px;
		height: 48px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--gradient-primary);
		border-radius: var(--radius-md);
		color: white;
		font-size: 1.25rem;
	}
	
	.form-actions {
		display: flex;
		gap: var(--space-md);
		justify-content: flex-end;
		margin-top: var(--space-xl);
	}
	
	.spin {
		animation: spin 1s linear infinite;
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
