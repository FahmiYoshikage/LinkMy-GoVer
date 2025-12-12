<script lang="ts">
	import { onMount } from 'svelte';
	
	const API_URL = 'http://localhost:3000';
	
	interface AdminUser {
		id: number;
		username: string;
		email: string;
		is_verified: boolean;
		is_active: boolean;
		is_admin: boolean;
		created_at: string;
		profile_count: number;
		total_clicks: number;
	}
	
	let users = $state<AdminUser[]>([]);
	let loading = $state(true);
	let search = $state('');
	let updating = $state<number | null>(null);
	
	onMount(() => loadUsers());
	
	async function loadUsers() {
		loading = true;
		const token = localStorage.getItem('access_token');
		const res = await fetch(`${API_URL}/api/v1/admin/users?search=${encodeURIComponent(search)}`, {
			headers: { Authorization: `Bearer ${token}` }
		});
		if (res.ok) {
			const data = await res.json();
			users = data.data || [];
		}
		loading = false;
	}
	
	async function updateUser(id: number, updates: Partial<AdminUser>) {
		updating = id;
		const token = localStorage.getItem('access_token');
		await fetch(`${API_URL}/api/v1/admin/users/${id}`, {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(updates)
		});
		// Update local state
		users = users.map(u => u.id === id ? { ...u, ...updates } : u);
		updating = null;
	}
	
	function formatDate(dateStr: string) {
		return new Date(dateStr).toLocaleDateString('id-ID', {
			year: 'numeric', month: 'short', day: 'numeric'
		});
	}
</script>

<svelte:head>
	<title>Users - Admin</title>
</svelte:head>

<div class="admin-users">
	<header class="page-header">
		<div class="header-left">
			<h1>Users</h1>
			<p>{users.length} users</p>
		</div>
		<div class="header-actions">
			<div class="search-box">
				<i class="bi bi-search"></i>
				<input 
					type="text" 
					placeholder="Search users..." 
					bind:value={search}
					onkeyup={(e) => e.key === 'Enter' && loadUsers()}
				/>
			</div>
		</div>
	</header>
	
	{#if loading}
		<div class="loading">Loading users...</div>
	{:else}
		<div class="table-container card">
			<table class="data-table">
				<thead>
					<tr>
						<th>User</th>
						<th>Profiles</th>
						<th>Clicks</th>
						<th>Joined</th>
						<th>Status</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each users as user}
						<tr class:inactive={!user.is_active}>
							<td>
								<div class="user-cell">
									<div class="avatar" class:admin={user.is_admin}>
										{user.is_admin ? '👑' : user.username.charAt(0).toUpperCase()}
									</div>
									<div class="user-info">
										<span class="username">
											{user.username}
											{#if user.is_verified}
												<i class="bi bi-patch-check-fill verified"></i>
											{/if}
										</span>
										<span class="email">{user.email}</span>
									</div>
								</div>
							</td>
							<td>{user.profile_count}</td>
							<td>{user.total_clicks.toLocaleString()}</td>
							<td>{formatDate(user.created_at)}</td>
							<td>
								<div class="status-badges">
									{#if user.is_admin}
										<span class="badge admin">Admin</span>
									{/if}
									{#if user.is_verified}
										<span class="badge verified">Verified</span>
									{:else}
										<span class="badge unverified">Unverified</span>
									{/if}
									{#if !user.is_active}
										<span class="badge banned">Banned</span>
									{/if}
								</div>
							</td>
							<td>
								<div class="actions">
									<button 
										class="btn-icon" 
										title={user.is_verified ? 'Remove verification' : 'Verify user'}
										onclick={() => updateUser(user.id, { is_verified: !user.is_verified })}
										disabled={updating === user.id}
									>
										<i class="bi {user.is_verified ? 'bi-patch-check-fill' : 'bi-patch-check'}"></i>
									</button>
									<button 
										class="btn-icon" 
										title={user.is_active ? 'Ban user' : 'Unban user'}
										onclick={() => updateUser(user.id, { is_active: !user.is_active })}
										disabled={updating === user.id}
									>
										<i class="bi {user.is_active ? 'bi-slash-circle' : 'bi-check-circle'}"></i>
									</button>
									<button 
										class="btn-icon" 
										title={user.is_admin ? 'Remove admin' : 'Make admin'}
										onclick={() => updateUser(user.id, { is_admin: !user.is_admin })}
										disabled={updating === user.id}
									>
										<i class="bi {user.is_admin ? 'bi-shield-slash' : 'bi-shield-check'}"></i>
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
	.admin-users {
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
	
	.search-box i {
		color: var(--color-text-muted);
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
		letter-spacing: 0.05em;
	}
	
	.data-table tr.inactive {
		opacity: 0.5;
	}
	
	.user-cell {
		display: flex;
		align-items: center;
		gap: var(--space-md);
	}
	
	.avatar {
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-bg-tertiary);
		border-radius: 50%;
		font-weight: 600;
	}
	
	.avatar.admin {
		background: linear-gradient(135deg, #fbbf24, #f59e0b);
	}
	
	.user-info {
		display: flex;
		flex-direction: column;
	}
	
	.username {
		font-weight: 500;
		display: flex;
		align-items: center;
		gap: var(--space-xs);
	}
	
	.verified {
		color: #3b82f6;
		font-size: 0.875rem;
	}
	
	.email {
		font-size: 0.75rem;
		color: var(--color-text-muted);
	}
	
	.status-badges {
		display: flex;
		gap: var(--space-xs);
		flex-wrap: wrap;
	}
	
	.badge {
		padding: 2px 8px;
		border-radius: 999px;
		font-size: 0.675rem;
		font-weight: 500;
		text-transform: uppercase;
	}
	
	.badge.admin { background: rgba(251, 191, 36, 0.2); color: #fbbf24; }
	.badge.verified { background: rgba(59, 130, 246, 0.2); color: #3b82f6; }
	.badge.unverified { background: rgba(107, 114, 128, 0.2); color: #6b7280; }
	.badge.banned { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
	
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
	
	.btn-icon:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
