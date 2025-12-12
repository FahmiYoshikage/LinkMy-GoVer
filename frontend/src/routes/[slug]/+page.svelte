<script lang="ts">
	import { onMount } from 'svelte';
	import { profiles, links as linksApi } from '$lib/api';
	import type { PublicProfile, Link } from '$lib/api';
	import ShareModal from '$lib/components/ShareModal.svelte';
	
	let { data } = $props();
	
	let profile = $state<PublicProfile | null>(null);
	let loading = $state(true);
	let error = $state('');
	let showShare = $state(false);
	
	onMount(async () => {
		const res = await profiles.getPublic(data.slug);
		if (res.error) {
			error = 'Profile not found';
		} else if (res.data) {
			profile = res.data;
		}
		loading = false;
	});
	
	async function handleClick(link: Link) {
		// Track click
		linksApi.trackClick(link.id);
		// Open URL
		window.open(link.url, '_blank');
	}
	
	function getBackground(): string {
		if (!profile?.theme) return 'var(--gradient-primary)';
		const t = profile.theme;
		
		if (t.boxed_enabled && t.boxed_outer_bg_value) {
			return t.boxed_outer_bg_value;
		}
		return t.bg_value || 'var(--gradient-primary)';
	}
	
	function getContainerBg(): string {
		if (!profile?.theme?.boxed_enabled) return 'transparent';
		return profile.theme.boxed_container_bg || '#ffffff';
	}
</script>

<svelte:head>
	{#if profile}
		<title>{profile.profile.title || profile.profile.name} - LinkMy</title>
		<meta name="description" content={profile.profile.bio || 'Check out my links!'} />
	{:else}
		<title>LinkMy</title>
	{/if}
</svelte:head>

{#if loading}
	<div class="loading-screen">
		<div class="loader"></div>
	</div>
{:else if error}
	<div class="error-screen">
		<div class="error-content">
			<i class="bi bi-emoji-frown"></i>
			<h1>404</h1>
			<p>Profile not found</p>
			<a href="/" class="btn btn-primary">Go Home</a>
		</div>
	</div>
{:else if profile}
	<div 
		class="profile-page" 
		style="background: {getBackground()}; --text-color: {profile.theme.text_color}; --btn-color: {profile.theme.button_color}; --font: '{profile.theme.font}', sans-serif;"
		class:boxed={profile.theme.boxed_enabled}
		class:animated={profile.theme.enable_animations}
		class:glass={profile.theme.enable_glass_effect}
	>
		<div 
			class="profile-container"
			class:container-boxed={profile.theme.boxed_enabled}
			style="
				{profile.theme.boxed_enabled ? `
					background: ${getContainerBg()};
					max-width: ${profile.theme.boxed_max_width}px;
					border-radius: ${profile.theme.boxed_radius}px;
				` : ''}
			"
		>
			<!-- Avatar -->
			<div class="profile-avatar">
				<img 
				src={profile.profile.avatar && profile.profile.avatar !== 'default-avatar.png' 
					? `/uploads/${profile.profile.avatar}` 
					: `https://ui-avatars.com/api/?name=${encodeURIComponent(profile.profile.name)}&background=667eea&color=fff&size=120`}
				alt={profile.profile.name}
			/>
				{#if profile.is_verified}
					<span class="verified-badge" title="Verified">
						<i class="bi bi-check-circle-fill"></i>
					</span>
				{/if}
			</div>
			
			<!-- Info -->
			<div class="profile-info">
				<h1>{profile.profile.title || profile.profile.name}</h1>
				{#if profile.profile.bio}
					<p>{profile.profile.bio}</p>
				{/if}
			</div>
			
			<!-- Links -->
			<div class="links-container">
				{#each profile.links as link, i}
					<button 
						class="link-button"
						class:pill={profile.theme.button_style === 'pill'}
						class:rounded={profile.theme.button_style === 'rounded'}
						class:square={profile.theme.button_style === 'square'}
						style="animation-delay: {i * 50}ms"
						onclick={() => handleClick(link)}
					>
						<i class="bi {link.icon}"></i>
						<span>{link.title}</span>
						<i class="bi bi-arrow-right link-arrow"></i>
					</button>
				{/each}
			</div>
			
			<!-- Footer -->
			<div class="profile-footer">
				<a href="/" class="powered-by">
					<i class="bi bi-diagram-3-fill"></i>
					<span>LinkMy</span>
				</a>
			</div>
		</div>
		
		<!-- Floating Share Button -->
		<button class="share-fab" onclick={() => showShare = true}>
			<i class="bi bi-share-fill"></i>
		</button>
	</div>
	
	<ShareModal bind:isOpen={showShare} slug={data.slug} profileName={profile.profile.name} />
{/if}

<style>
	.loading-screen, .error-screen {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-bg);
	}
	
	.loader {
		width: 48px;
		height: 48px;
		border: 4px solid var(--color-border);
		border-top-color: var(--color-primary);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
	
	.error-content {
		text-align: center;
	}
	
	.error-content i {
		font-size: 4rem;
		color: var(--color-text-muted);
		margin-bottom: var(--space-lg);
	}
	
	.error-content h1 {
		font-size: 4rem;
		margin-bottom: var(--space-sm);
	}
	
	.error-content p {
		color: var(--color-text-secondary);
		margin-bottom: var(--space-xl);
	}
	
	/* Profile Page */
	.profile-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-lg);
		font-family: var(--font);
		color: var(--text-color);
	}
	
	.profile-container {
		width: 100%;
		max-width: 680px;
		text-align: center;
		padding: var(--space-2xl) var(--space-lg);
	}
	
	.container-boxed {
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
	}
	
	/* Avatar */
	.profile-avatar {
		position: relative;
		display: inline-block;
		margin-bottom: var(--space-lg);
	}
	
	.profile-avatar img {
		width: 120px;
		height: 120px;
		border-radius: 50%;
		object-fit: cover;
		border: 4px solid rgba(255, 255, 255, 0.2);
	}
	
	.verified-badge {
		position: absolute;
		bottom: 5px;
		right: 5px;
		width: 28px;
		height: 28px;
		background: var(--color-primary);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
		font-size: 0.875rem;
	}
	
	/* Info */
	.profile-info h1 {
		font-size: 1.75rem;
		font-weight: 700;
		margin-bottom: var(--space-sm);
	}
	
	.profile-info p {
		opacity: 0.8;
		margin-bottom: var(--space-xl);
	}
	
	/* Links */
	.links-container {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		margin-bottom: var(--space-2xl);
	}
	
	.link-button {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		width: 100%;
		padding: var(--space-lg);
		background: var(--btn-color);
		color: white;
		border: none;
		cursor: pointer;
		font-size: 1rem;
		font-weight: 500;
		transition: all 0.2s ease;
	}
	
	.link-button.animated {
		animation: slideUp 0.5s ease backwards;
	}
	
	.link-button.pill {
		border-radius: 50px;
	}
	
	.link-button.rounded {
		border-radius: var(--radius-lg);
	}
	
	.link-button.square {
		border-radius: 0;
	}
	
	.link-button:hover {
		transform: scale(1.02);
		box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
	}
	
	.link-button span {
		flex: 1;
		text-align: left;
	}
	
	.link-arrow {
		opacity: 0;
		transition: opacity 0.2s;
	}
	
	.link-button:hover .link-arrow {
		opacity: 1;
	}
	
	/* Glass Effect */
	.glass .link-button {
		background: rgba(255, 255, 255, 0.15);
		backdrop-filter: blur(10px);
		-webkit-backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.2);
	}
	
	/* Footer */
	.profile-footer {
		margin-top: var(--space-2xl);
	}
	
	.powered-by {
		display: inline-flex;
		align-items: center;
		gap: var(--space-xs);
		font-size: 0.875rem;
		color: inherit;
		opacity: 0.6;
		transition: opacity 0.2s;
	}
	
	.powered-by:hover {
		opacity: 1;
	}
	
	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
	
	/* Floating Share Button */
	.share-fab {
		position: fixed;
		bottom: 2rem;
		right: 2rem;
		width: 56px;
		height: 56px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border: none;
		border-radius: 50%;
		color: white;
		font-size: 1.25rem;
		cursor: pointer;
		box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
		transition: transform 0.2s, box-shadow 0.2s;
		z-index: 100;
	}
	
	.share-fab:hover {
		transform: scale(1.1);
		box-shadow: 0 6px 30px rgba(102, 126, 234, 0.6);
	}
	
	.share-fab:active {
		transform: scale(0.95);
	}
</style>
