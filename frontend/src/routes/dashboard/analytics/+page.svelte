<script lang="ts">
	import { onMount } from 'svelte';
	import { profiles, type Profile } from '$lib/api';
	
	let userProfiles = $state<Profile[]>([]);
	let selectedProfileId = $state<number | null>(null);
	let analytics = $state<any>(null);
	let loading = $state(true);
	let days = $state(30);
	
	const API_URL = 'http://localhost:3000';
	
	onMount(async () => {
		const res = await profiles.getAll();
		if (res.data && res.data.length > 0) {
			userProfiles = res.data;
			selectedProfileId = res.data[0].id;
			await loadAnalytics();
		}
		loading = false;
	});
	
	async function loadAnalytics() {
		if (!selectedProfileId) return;
		const token = localStorage.getItem('access_token');
		const res = await fetch(`${API_URL}/api/v1/profiles/${selectedProfileId}/analytics?days=${days}`, {
			headers: { Authorization: `Bearer ${token}` }
		});
		if (res.ok) {
			const data = await res.json();
			analytics = data.data;
		}
	}
	
	async function handleProfileChange(e: Event) {
		const select = e.target as HTMLSelectElement;
		selectedProfileId = parseInt(select.value);
		await loadAnalytics();
	}
	
	async function handleDaysChange(newDays: number) {
		days = newDays;
		await loadAnalytics();
	}
</script>

<svelte:head>
	<title>Analytics - LinkMy</title>
</svelte:head>

<div class="analytics-page">
	<header class="page-header">
		<div class="header-left">
			<h1>Analytics</h1>
			<p>Track your link performance</p>
		</div>
		<div class="header-actions">
			{#if userProfiles.length > 1}
				<select class="profile-select" onchange={handleProfileChange} value={selectedProfileId}>
					{#each userProfiles as profile}
						<option value={profile.id}>{profile.name}</option>
					{/each}
				</select>
			{/if}
			<div class="days-filter">
				<button class:active={days === 7} onclick={() => handleDaysChange(7)}>7d</button>
				<button class:active={days === 30} onclick={() => handleDaysChange(30)}>30d</button>
				<button class:active={days === 90} onclick={() => handleDaysChange(90)}>90d</button>
			</div>
		</div>
	</header>
	
	{#if loading}
		<div class="loading-state">
			<div class="loader"></div>
		</div>
	{:else if analytics}
		<!-- Stats Cards -->
		<div class="stats-grid">
			<div class="stat-card card">
				<div class="stat-icon">
					<i class="bi bi-cursor-fill"></i>
				</div>
				<div class="stat-info">
					<span class="stat-value">{analytics.total_clicks || 0}</span>
					<span class="stat-label">Total Clicks</span>
				</div>
			</div>
			
			<div class="stat-card card">
				<div class="stat-icon views">
					<i class="bi bi-eye-fill"></i>
				</div>
				<div class="stat-info">
					<span class="stat-value">{analytics.clicks_by_day?.length || 0}</span>
					<span class="stat-label">Active Days</span>
				</div>
			</div>
			
			<div class="stat-card card">
				<div class="stat-icon links">
					<i class="bi bi-link-45deg"></i>
				</div>
				<div class="stat-info">
					<span class="stat-value">{analytics.clicks_by_link?.length || 0}</span>
					<span class="stat-label">Active Links</span>
				</div>
			</div>
		</div>
		
		<div class="analytics-grid">
			<!-- Clicks Chart -->
			<section class="chart-section card">
				<h2>Clicks Over Time</h2>
				{#if analytics.clicks_by_day?.length > 0}
					<div class="chart-area">
						<div class="bar-chart">
							{#each analytics.clicks_by_day as day}
								<div class="bar-column">
									<div 
										class="bar" 
										style="height: {Math.max(5, (day.clicks / Math.max(...analytics.clicks_by_day.map((d: any) => d.clicks))) * 100)}%"
										title="{day.date}: {day.clicks} clicks"
									></div>
									<span class="bar-label">{day.date.slice(5)}</span>
								</div>
							{/each}
						</div>
					</div>
				{:else}
					<div class="empty-chart">
						<i class="bi bi-graph-up"></i>
						<p>No click data yet</p>
					</div>
				{/if}
			</section>
			
			<!-- Top Links -->
			<section class="list-section card">
				<h2>Top Links</h2>
				{#if analytics.clicks_by_link?.length > 0}
					<div class="top-list">
						{#each analytics.clicks_by_link as link, i}
							<div class="list-item">
								<span class="rank">#{i + 1}</span>
								<span class="title">{link.title}</span>
								<span class="clicks">{link.clicks}</span>
							</div>
						{/each}
					</div>
				{:else}
					<div class="empty-list">
						<p>No link clicks yet</p>
					</div>
				{/if}
			</section>
			
			<!-- Top Countries -->
			<section class="list-section card">
				<h2>Top Countries</h2>
				{#if analytics.clicks_by_country?.length > 0}
					<div class="top-list">
						{#each analytics.clicks_by_country as country, i}
							<div class="list-item">
								<span class="rank">#{i + 1}</span>
								<span class="title">{country.country}</span>
								<span class="clicks">{country.clicks}</span>
							</div>
						{/each}
					</div>
				{:else}
					<div class="empty-list">
						<p>No country data yet</p>
					</div>
				{/if}
			</section>
			
			<!-- Top Referrers -->
			<section class="list-section card">
				<h2>Top Referrers</h2>
				{#if analytics.top_referrers?.length > 0}
					<div class="top-list">
						{#each analytics.top_referrers as ref, i}
							<div class="list-item">
								<span class="rank">#{i + 1}</span>
								<span class="title">{ref.referrer || 'Direct'}</span>
								<span class="clicks">{ref.clicks}</span>
							</div>
						{/each}
					</div>
				{:else}
					<div class="empty-list">
						<p>No referrer data yet</p>
					</div>
				{/if}
			</section>
		</div>
	{/if}
</div>

<style>
	.analytics-page {
		max-width: 1200px;
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
	}
	
	.days-filter {
		display: flex;
		background: var(--color-bg-tertiary);
		border-radius: var(--radius-lg);
		overflow: hidden;
	}
	
	.days-filter button {
		padding: var(--space-md) var(--space-lg);
		color: var(--color-text-secondary);
		transition: all var(--transition-fast);
	}
	
	.days-filter button.active {
		background: var(--color-primary);
		color: white;
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
	
	.stat-icon.views {
		background: linear-gradient(135deg, #10b981 0%, #059669 100%);
	}
	
	.stat-icon.links {
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
	
	/* Analytics Grid */
	.analytics-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--space-lg);
	}
	
	.chart-section {
		grid-column: 1 / -1;
		padding: var(--space-xl);
	}
	
	.chart-section h2,
	.list-section h2 {
		font-size: 1rem;
		margin-bottom: var(--space-lg);
		color: var(--color-text-secondary);
	}
	
	.list-section {
		padding: var(--space-xl);
	}
	
	/* Bar Chart */
	.chart-area {
		height: 200px;
	}
	
	.bar-chart {
		display: flex;
		align-items: flex-end;
		gap: var(--space-xs);
		height: 100%;
	}
	
	.bar-column {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		height: 100%;
	}
	
	.bar {
		width: 100%;
		max-width: 30px;
		background: var(--gradient-primary);
		border-radius: var(--radius-sm) var(--radius-sm) 0 0;
		margin-top: auto;
	}
	
	.bar-label {
		font-size: 0.625rem;
		color: var(--color-text-muted);
		margin-top: var(--space-xs);
	}
	
	.empty-chart {
		height: 200px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		color: var(--color-text-muted);
	}
	
	.empty-chart i {
		font-size: 2rem;
		margin-bottom: var(--space-md);
	}
	
	/* Top Lists */
	.top-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
	}
	
	.list-item {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		padding: var(--space-sm) 0;
		border-bottom: 1px solid var(--color-border);
	}
	
	.list-item:last-child {
		border-bottom: none;
	}
	
	.rank {
		font-size: 0.75rem;
		color: var(--color-text-muted);
		width: 24px;
	}
	
	.title {
		flex: 1;
		font-size: 0.875rem;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.clicks {
		font-weight: 600;
		color: var(--color-primary);
	}
	
	.empty-list {
		padding: var(--space-xl);
		text-align: center;
		color: var(--color-text-muted);
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
	
	@media (max-width: 768px) {
		.analytics-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
