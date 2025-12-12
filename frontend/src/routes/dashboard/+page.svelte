<script lang="ts">
	import { onMount } from 'svelte';
	import { profiles } from '$lib/api';
	import type { Profile } from '$lib/api';
	
	let userProfiles = $state<Profile[]>([]);
	let loading = $state(true);
	
	onMount(async () => {
		const res = await profiles.getAll();
		if (res.data) {
			userProfiles = res.data;
		}
		loading = false;
	});
</script>

<svelte:head>
	<title>Dashboard - LinkMy</title>
</svelte:head>

<div class="dashboard-page">
	<header class="page-header">
		<h1>Dashboard</h1>
		<p>Welcome back! Manage your link pages.</p>
	</header>
	
	<!-- Stats Cards -->
	<div class="stats-grid">
		<div class="stat-card card">
			<div class="stat-icon">
				<i class="bi bi-layers-fill"></i>
			</div>
			<div class="stat-content">
				<span class="stat-value">{userProfiles.length}</span>
				<span class="stat-label">Profiles</span>
			</div>
		</div>
		
		<div class="stat-card card">
			<div class="stat-icon links">
				<i class="bi bi-link-45deg"></i>
			</div>
			<div class="stat-content">
				<span class="stat-value">
					{userProfiles.reduce((acc, p) => acc + (p.link_count || 0), 0)}
				</span>
				<span class="stat-label">Total Links</span>
			</div>
		</div>
		
		<div class="stat-card card">
			<div class="stat-icon clicks">
				<i class="bi bi-cursor-fill"></i>
			</div>
			<div class="stat-content">
				<span class="stat-value">
					{userProfiles.reduce((acc, p) => acc + (p.total_clicks || 0), 0)}
				</span>
				<span class="stat-label">Total Clicks</span>
			</div>
		</div>
	</div>
	
	<!-- Profiles Section -->
	<section class="profiles-section">
		<div class="section-header">
			<h2>Your Profiles</h2>
			<a href="/dashboard/profiles/new" class="btn btn-primary">
				<i class="bi bi-plus-lg"></i>
				New Profile
			</a>
		</div>
		
		{#if loading}
			<div class="loading-state">
				<div class="loader"></div>
			</div>
		{:else if userProfiles.length === 0}
			<div class="empty-state card">
				<i class="bi bi-folder-x"></i>
				<h3>No profiles yet</h3>
				<p>Create your first profile to get started</p>
				<a href="/dashboard/profiles/new" class="btn btn-primary">
					<i class="bi bi-plus-lg"></i>
					Create Profile
				</a>
			</div>
		{:else}
			<div class="profiles-grid">
				{#each userProfiles as profile}
					<div class="profile-card card">
						<div class="profile-avatar">
							<img 
							src={profile.avatar && profile.avatar !== 'default-avatar.png' 
								? `/uploads/${profile.avatar}` 
								: `https://ui-avatars.com/api/?name=${encodeURIComponent(profile.name)}&background=667eea&color=fff&size=60`}
							alt={profile.name}
						/>
						</div>
						<div class="profile-info">
							<h3>{profile.title || profile.name}</h3>
							<a href="/{profile.slug}" target="_blank" class="profile-url">
								linkmy.deepkernel.site/{profile.slug}
								<i class="bi bi-box-arrow-up-right"></i>
							</a>
						</div>
						<div class="profile-stats">
							<div class="profile-stat">
								<span class="value">{profile.link_count || 0}</span>
								<span class="label">Links</span>
							</div>
							<div class="profile-stat">
								<span class="value">{profile.total_clicks || 0}</span>
								<span class="label">Clicks</span>
							</div>
						</div>
						<div class="profile-actions">
							<a href="/dashboard/links?profile={profile.id}" class="btn btn-secondary btn-sm">
								<i class="bi bi-pencil"></i>
								Edit
							</a>
							<button class="btn btn-ghost btn-sm" onclick={() => navigator.clipboard.writeText(`https://linkmy.deepkernel.site/${profile.slug}`)}>
								<i class="bi bi-clipboard"></i>
							</button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</section>
</div>

<style>
	.dashboard-page {
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
	
	/* Stats */
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: var(--space-lg);
		margin-bottom: var(--space-2xl);
	}
	
	.stat-card {
		display: flex;
		align-items: center;
		gap: var(--space-lg);
		padding: var(--space-xl);
	}
	
	.stat-icon {
		width: 56px;
		height: 56px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--gradient-primary);
		border-radius: var(--radius-lg);
		font-size: 1.5rem;
		color: white;
	}
	
	.stat-icon.links {
		background: linear-gradient(135deg, #10b981 0%, #059669 100%);
	}
	
	.stat-icon.clicks {
		background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
	}
	
	.stat-value {
		display: block;
		font-size: 1.75rem;
		font-weight: 700;
	}
	
	.stat-label {
		color: var(--color-text-secondary);
		font-size: 0.875rem;
	}
	
	/* Profiles Section */
	.section-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-lg);
	}
	
	.section-header h2 {
		font-size: 1.25rem;
	}
	
	.profiles-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
		gap: var(--space-lg);
	}
	
	.profile-card {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
	}
	
	.profile-avatar img {
		width: 60px;
		height: 60px;
		border-radius: var(--radius-lg);
		object-fit: cover;
	}
	
	.profile-info h3 {
		font-size: 1.125rem;
		margin-bottom: var(--space-xs);
	}
	
	.profile-url {
		font-size: 0.875rem;
		color: var(--color-text-muted);
		display: inline-flex;
		align-items: center;
		gap: var(--space-xs);
	}
	
	.profile-url:hover {
		color: var(--color-primary);
	}
	
	.profile-stats {
		display: flex;
		gap: var(--space-xl);
		padding: var(--space-md) 0;
		border-top: 1px solid var(--color-border);
		border-bottom: 1px solid var(--color-border);
	}
	
	.profile-stat {
		display: flex;
		flex-direction: column;
	}
	
	.profile-stat .value {
		font-weight: 600;
	}
	
	.profile-stat .label {
		font-size: 0.75rem;
		color: var(--color-text-muted);
	}
	
	.profile-actions {
		display: flex;
		gap: var(--space-sm);
	}
	
	.btn-sm {
		padding: var(--space-sm) var(--space-md);
		font-size: 0.875rem;
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
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
