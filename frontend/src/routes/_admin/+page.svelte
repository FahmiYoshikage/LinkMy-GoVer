<script lang="ts">
	import { onMount } from 'svelte';
	
	const API_URL = 'http://localhost:3000';
	
	interface Stats {
		total_users: number;
		verified_users: number;
		total_profiles: number;
		total_links: number;
		total_clicks: number;
		new_users_week: number;
		active_users_week: number;
	}
	
	let stats = $state<Stats | null>(null);
	let loading = $state(true);
	
	onMount(async () => {
		const token = localStorage.getItem('access_token');
		const res = await fetch(`${API_URL}/api/v1/admin/stats`, {
			headers: { Authorization: `Bearer ${token}` }
		});
		if (res.ok) {
			const data = await res.json();
			stats = data.data;
		}
		loading = false;
	});
</script>

<svelte:head>
	<title>Admin Dashboard - LinkMy</title>
</svelte:head>

<div class="admin-dashboard">
	<header class="page-header">
		<h1>Dashboard</h1>
		<p>Platform overview</p>
	</header>
	
	{#if loading}
		<div class="loading">Loading...</div>
	{:else if stats}
		<div class="stats-grid">
			<div class="stat-card">
				<div class="stat-icon users">
					<i class="bi bi-people-fill"></i>
				</div>
				<div class="stat-info">
					<span class="stat-value">{stats.total_users}</span>
					<span class="stat-label">Total Users</span>
				</div>
			</div>
			
			<div class="stat-card">
				<div class="stat-icon verified">
					<i class="bi bi-patch-check-fill"></i>
				</div>
				<div class="stat-info">
					<span class="stat-value">{stats.verified_users}</span>
					<span class="stat-label">Verified Users</span>
				</div>
			</div>
			
			<div class="stat-card">
				<div class="stat-icon profiles">
					<i class="bi bi-layers-fill"></i>
				</div>
				<div class="stat-info">
					<span class="stat-value">{stats.total_profiles}</span>
					<span class="stat-label">Profiles</span>
				</div>
			</div>
			
			<div class="stat-card">
				<div class="stat-icon links">
					<i class="bi bi-link-45deg"></i>
				</div>
				<div class="stat-info">
					<span class="stat-value">{stats.total_links}</span>
					<span class="stat-label">Links</span>
				</div>
			</div>
			
			<div class="stat-card">
				<div class="stat-icon clicks">
					<i class="bi bi-cursor-fill"></i>
				</div>
				<div class="stat-info">
					<span class="stat-value">{stats.total_clicks.toLocaleString()}</span>
					<span class="stat-label">Total Clicks</span>
				</div>
			</div>
			
			<div class="stat-card">
				<div class="stat-icon new">
					<i class="bi bi-person-plus-fill"></i>
				</div>
				<div class="stat-info">
					<span class="stat-value">{stats.new_users_week}</span>
					<span class="stat-label">New Users (7d)</span>
				</div>
			</div>
		</div>
		
		<div class="quick-actions">
			<h2>Quick Actions</h2>
			<div class="action-grid">
				<a href="/admin/users" class="action-card">
					<i class="bi bi-people"></i>
					<span>Manage Users</span>
				</a>
				<a href="/admin/profiles" class="action-card">
					<i class="bi bi-layers"></i>
					<span>Moderate Profiles</span>
				</a>
			</div>
		</div>
	{/if}
</div>

<style>
	.admin-dashboard {
		max-width: 1200px;
	}
	
	.page-header {
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
	
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: var(--space-lg);
		margin-bottom: var(--space-3xl);
	}
	
	.stat-card {
		display: flex;
		align-items: center;
		gap: var(--space-lg);
		padding: var(--space-xl);
		background: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
	}
	
	.stat-icon {
		width: 56px;
		height: 56px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: var(--radius-lg);
		font-size: 1.5rem;
	}
	
	.stat-icon.users { background: rgba(102, 126, 234, 0.2); color: #667eea; }
	.stat-icon.verified { background: rgba(16, 185, 129, 0.2); color: #10b981; }
	.stat-icon.profiles { background: rgba(139, 92, 246, 0.2); color: #8b5cf6; }
	.stat-icon.links { background: rgba(236, 72, 153, 0.2); color: #ec4899; }
	.stat-icon.clicks { background: rgba(251, 191, 36, 0.2); color: #fbbf24; }
	.stat-icon.new { background: rgba(6, 182, 212, 0.2); color: #06b6d4; }
	
	.stat-info {
		display: flex;
		flex-direction: column;
	}
	
	.stat-value {
		font-size: 1.75rem;
		font-weight: 700;
	}
	
	.stat-label {
		font-size: 0.875rem;
		color: var(--color-text-secondary);
	}
	
	.quick-actions h2 {
		font-size: 1.25rem;
		margin-bottom: var(--space-lg);
	}
	
	.action-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: var(--space-lg);
	}
	
	.action-card {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-md);
		padding: var(--space-xl);
		background: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
		color: var(--color-text);
		transition: all var(--transition-fast);
	}
	
	.action-card:hover {
		border-color: var(--color-primary);
		transform: translateY(-2px);
	}
	
	.action-card i {
		font-size: 2rem;
		color: var(--color-primary);
	}
</style>
