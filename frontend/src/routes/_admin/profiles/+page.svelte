<script lang="ts">
	import { onMount } from 'svelte';
	
	const API_URL = 'http://localhost:3000';
	
	interface AdminProfile {
		id: number;
		user_id: number;
		username: string;
		slug: string;
		name: string;
		is_active: boolean;
		link_count: number;
		total_clicks: number;
		created_at: string;
	}
	
	let profiles = $state<AdminProfile[]>([]);
	let loading = $state(true);
	let search = $state('');
	let updating = $state<number | null>(null);
	
	onMount(() => loadProfiles());
	
	async function loadProfiles() {
		loading = true;
		const token = localStorage.getItem('access_token');
		const res = await fetch(`${API_URL}/api/v1/admin/profiles?search=${encodeURIComponent(search)}`, {
			headers: { Authorization: `Bearer ${token}` }
		});
		if (res.ok) {
			const data = await res.json();
			profiles = data.data || [];
		}
		loading = false;
	}
	
	async function toggleActive(id: number, current: boolean) {
		updating = id;
		const token = localStorage.getItem('access_token');
		await fetch(`${API_URL}/api/v1/admin/profiles/${id}`, {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ is_active: !current })
		});
		profiles = profiles.map(p => p.id === id ? { ...p, is_active: !current } : p);
		updating = null;
	}
	
	function formatDate(dateStr: string) {
		return new Date(dateStr).toLocaleDateString('id-ID', {
			year: 'numeric', month: 'short', day: 'numeric'
		});
	}
</script>

<svelte:head>
	<title>Profiles - Admin</title>
</svelte:head>

<div class="admin-profiles">
	<header class="page-header">
		<div class="header-left">
			<h1>Profiles</h1>
			<p>{profiles.length} profiles</p>
		</div>
		<div class="header-actions">
			<div class="search-box">
				<i class="bi bi-search"></i>
				<input 
					type="text" 
					placeholder="Search profiles..." 
					bind:value={search}
					onkeyup={(e) => e.key === 'Enter' && loadProfiles()}
				/>
			</div>
		</div>
	</header>
	
	{#if loading}
		<div class="loading">Loading profiles...</div>
	{:else}
		<div class="table-container card">
			<table class="data-table">
				<thead>
					<tr>
						<th>Profile</th>
						<th>Owner</th>
						<th>Links</th>
						<th>Clicks</th>
						<th>Created</th>
						<th>Status</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each profiles as profile}
						<tr class:inactive={!profile.is_active}>
							<td>
								<div class="profile-cell">
									<span class="profile-name">{profile.name}</span>
									<a href="/{profile.slug}" target="_blank" class="profile-slug">
										/{profile.slug}
									</a>
								</div>
							</td>
							<td>
								<span class="username">@{profile.username}</span>
							</td>
							<td>{profile.link_count}</td>
							<td>{profile.total_clicks.toLocaleString()}</td>
							<td>{formatDate(profile.created_at)}</td>
							<td>
								{#if profile.is_active}
									<span class="badge active">Active</span>
								{:else}
									<span class="badge hidden">Hidden</span>
								{/if}
							</td>
							<td>
								<div class="actions">
									<a 
										href="/{profile.slug}" 
										target="_blank"
										class="btn-icon" 
										title="View profile"
									>
										<i class="bi bi-eye"></i>
									</a>
									<button 
										class="btn-icon" 
										title={profile.is_active ? 'Hide profile' : 'Show profile'}
										onclick={() => toggleActive(profile.id, profile.is_active)}
										disabled={updating === profile.id}
									>
										<i class="bi {profile.is_active ? 'bi-eye-slash' : 'bi-eye'}"></i>
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<style>
	.admin-profiles {
		max-width: 1200px;
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
	
	.search-box {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		padding: var(--space-sm) var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
	}
	
	.search-box input {
		background: none;
		border: none;
		color: var(--color-text);
		width: 200px;
	}
	
	.loading {
		text-align: center;
		padding: var(--space-3xl);
		color: var(--color-text-muted);
	}
	
	.table-container {
		overflow-x: auto;
	}
	
	.data-table {
		width: 100%;
		border-collapse: collapse;
	}
	
	.data-table th,
	.data-table td {
		padding: var(--space-md);
		text-align: left;
		border-bottom: 1px solid var(--color-border);
	}
	
	.data-table th {
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-text-muted);
		text-transform: uppercase;
	}
	
	.data-table tr.inactive {
		opacity: 0.5;
	}
	
	.profile-cell {
		display: flex;
		flex-direction: column;
	}
	
	.profile-name {
		font-weight: 500;
	}
	
	.profile-slug {
		font-size: 0.75rem;
		color: var(--color-primary);
	}
	
	.username {
		color: var(--color-text-secondary);
		font-size: 0.875rem;
	}
	
	.badge {
		padding: 2px 8px;
		border-radius: 999px;
		font-size: 0.675rem;
		font-weight: 500;
		text-transform: uppercase;
	}
	
	.badge.active { background: rgba(16, 185, 129, 0.2); color: #10b981; }
	.badge.hidden { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
	
	.actions {
		display: flex;
		gap: var(--space-xs);
	}
	
	.btn-icon {
		width: 32px;
		height: 32px;
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
</style>
