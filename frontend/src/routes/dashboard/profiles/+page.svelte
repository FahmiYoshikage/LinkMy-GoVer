<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { profiles, type Profile } from '$lib/api';
	
	let userProfiles = $state<Profile[]>([]);
	let loading = $state(true);
	let maxFreeProfiles = 2;
	
	onMount(async () => {
		const res = await profiles.getAll();
		if (res.data) {
			userProfiles = res.data;
		}
		loading = false;
	});
	
	async function deleteProfile(id: number) {
		if (!confirm('Are you sure you want to delete this profile? All links will be lost.')) return;
		
		await profiles.delete(id);
		userProfiles = userProfiles.filter(p => p.id !== id);
	}
	
	function copyLink(slug: string) {
		navigator.clipboard.writeText(`https://linkmy.deepkernel.site/${slug}`);
	}
</script>

<svelte:head>
	<title>Profiles - LinkMy</title>
</svelte:head>

<div class="profiles-page">
	<header class="page-header">
		<div class="header-left">
			<h1>Profiles</h1>
			<p>Manage your link pages ({userProfiles.length}/{maxFreeProfiles} free)</p>
		</div>
		<div class="header-actions">
			{#if userProfiles.length < maxFreeProfiles}
				<a href="/dashboard/profiles/new" class="btn btn-primary">
					<i class="bi bi-plus-lg"></i>
					New Profile
				</a>
			{:else}
				<button class="btn btn-secondary" disabled title="Upgrade to Pro for more profiles">
					<i class="bi bi-lock"></i>
					Upgrade for More
				</button>
			{/if}
		</div>
	</header>
	
	{#if loading}
		<div class="loading-state">
			<div class="loader"></div>
		</div>
	{:else if userProfiles.length === 0}
		<div class="empty-state card">
			<i class="bi bi-layers"></i>
			<h3>No profiles yet</h3>
			<p>Create your first profile to get started</p>
			<a href="/dashboard/profiles/new" class="btn btn-primary">
				<i class="bi bi-plus-lg"></i>
				Create Profile
			</a>
		</div>
	{:else}
		<div class="profiles-list">
			{#each userProfiles as profile}
				<div class="profile-card card">
					<div class="profile-avatar">
						<img
							src={profile.avatar && profile.avatar !== 'default-avatar.png' 
								? `/uploads/${profile.avatar}` 
								: `https://ui-avatars.com/api/?name=${encodeURIComponent(profile.name)}&background=667eea&color=fff&size=80`}
							alt={profile.name}
						/>
					</div>
					
					<div class="profile-info">
						<h3>{profile.name}</h3>
						<p class="profile-title">{profile.title || 'No title'}</p>
						<a href="/{profile.slug}" target="_blank" class="profile-url">
							linkmy.deepkernel.site/{profile.slug}
							<i class="bi bi-box-arrow-up-right"></i>
						</a>
					</div>
					
					<div class="profile-stats">
						<div class="stat">
							<span class="value">{profile.link_count || 0}</span>
							<span class="label">Links</span>
						</div>
						<div class="stat">
							<span class="value">{profile.total_clicks || 0}</span>
							<span class="label">Clicks</span>
						</div>
					</div>
					
					<div class="profile-actions">
						<a href="/dashboard/profiles/edit?id={profile.id}" class="btn btn-secondary btn-sm">
							<i class="bi bi-pencil"></i>
							Edit
						</a>
						<a href="/dashboard/links?profile={profile.id}" class="btn btn-secondary btn-sm">
							<i class="bi bi-link-45deg"></i>
							Links
						</a>
						<a href="/dashboard/appearance?profile={profile.id}" class="btn btn-secondary btn-sm">
							<i class="bi bi-palette"></i>
							Style
						</a>
						<button class="btn-icon" onclick={() => copyLink(profile.slug)} title="Copy Link">
							<i class="bi bi-clipboard"></i>
						</button>
						<button class="btn-icon danger" onclick={() => deleteProfile(profile.id)} title="Delete">
							<i class="bi bi-trash"></i>
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
	
	{#if userProfiles.length >= maxFreeProfiles}
		<div class="upgrade-banner card">
			<i class="bi bi-rocket-takeoff"></i>
			<div class="upgrade-content">
				<h3>Want more profiles?</h3>
				<p>Upgrade to Pro for unlimited profiles and more features - just $2/month</p>
			</div>
			<a href="/pricing" class="btn btn-primary">
				<i class="bi bi-stars"></i>
				Upgrade to Pro
			</a>
		</div>
	{/if}
</div>

<style>
	.profiles-page {
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
	
	.profiles-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}
	
	.profile-card {
		display: flex;
		align-items: center;
		gap: var(--space-lg);
		padding: var(--space-xl);
		flex-wrap: wrap;
	}
	
	.profile-avatar img {
		width: 80px;
		height: 80px;
		border-radius: var(--radius-lg);
		object-fit: cover;
	}
	
	.profile-info {
		flex: 1;
		min-width: 200px;
	}
	
	.profile-info h3 {
		font-size: 1.25rem;
		margin-bottom: var(--space-xs);
	}
	
	.profile-title {
		color: var(--color-text-secondary);
		margin-bottom: var(--space-sm);
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
	}
	
	.stat {
		text-align: center;
	}
	
	.stat .value {
		display: block;
		font-size: 1.5rem;
		font-weight: 600;
	}
	
	.stat .label {
		font-size: 0.75rem;
		color: var(--color-text-muted);
	}
	
	.profile-actions {
		display: flex;
		gap: var(--space-sm);
		flex-wrap: wrap;
	}
	
	.btn-sm {
		padding: var(--space-sm) var(--space-md);
		font-size: 0.875rem;
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
	
	/* Upgrade Banner */
	.upgrade-banner {
		display: flex;
		align-items: center;
		gap: var(--space-xl);
		padding: var(--space-xl);
		margin-top: var(--space-2xl);
		background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
		border-color: rgba(102, 126, 234, 0.3);
	}
	
	.upgrade-banner > i {
		font-size: 2rem;
		color: var(--color-primary);
	}
	
	.upgrade-content {
		flex: 1;
	}
	
	.upgrade-content h3 {
		margin-bottom: var(--space-xs);
	}
	
	.upgrade-content p {
		color: var(--color-text-secondary);
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
