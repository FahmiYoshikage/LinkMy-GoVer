<script lang="ts">
	import { goto } from '$app/navigation';
	import { profiles } from '$lib/api';
	
	let name = $state('');
	let slug = $state('');
	let title = $state('');
	let bio = $state('');
	let error = $state('');
	let loading = $state(false);
	
	// Auto-generate slug from name
	function generateSlug() {
		slug = name.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/^-|-$/g, '');
	}
	
	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;
		
		const res = await profiles.create({
			name,
			slug,
			title: title || undefined,
			bio: bio || undefined
		});
		
		if (res.error) {
			error = res.error;
			loading = false;
			return;
		}
		
		goto('/dashboard');
	}
</script>

<svelte:head>
	<title>New Profile - LinkMy</title>
</svelte:head>

<div class="new-profile-page">
	<header class="page-header">
		<a href="/dashboard" class="back-link">
			<i class="bi bi-arrow-left"></i>
			Back
		</a>
		<h1>Create New Profile</h1>
		<p>Set up a new link page</p>
	</header>
	
	<form class="profile-form card" onsubmit={handleSubmit}>
		{#if error}
			<div class="alert alert-error">
				<i class="bi bi-exclamation-circle"></i>
				{error}
			</div>
		{/if}
		
		<div class="form-group">
			<label for="name">Profile Name *</label>
			<input
				type="text"
				id="name"
				bind:value={name}
				oninput={generateSlug}
				placeholder="My Portfolio"
				required
			/>
		</div>
		
		<div class="form-group">
			<label for="slug">URL Slug *</label>
			<div class="slug-input">
				<span class="slug-prefix">linkmy.deepkernel.site/</span>
				<input
					type="text"
					id="slug"
					bind:value={slug}
					placeholder="my-portfolio"
					pattern="[a-z0-9-]+"
					required
				/>
			</div>
			<span class="input-hint">Only lowercase letters, numbers, and dashes</span>
		</div>
		
		<div class="form-group">
			<label for="title">Display Title</label>
			<input
				type="text"
				id="title"
				bind:value={title}
				placeholder="Web Developer & Designer"
			/>
		</div>
		
		<div class="form-group">
			<label for="bio">Bio</label>
			<textarea
				id="bio"
				bind:value={bio}
				placeholder="Tell visitors about yourself..."
				rows="3"
			></textarea>
		</div>
		
		<div class="form-actions">
			<a href="/dashboard" class="btn btn-secondary">Cancel</a>
			<button type="submit" class="btn btn-primary" disabled={loading}>
				{#if loading}
					<i class="bi bi-arrow-repeat spin"></i>
					Creating...
				{:else}
					<i class="bi bi-plus-lg"></i>
					Create Profile
				{/if}
			</button>
		</div>
	</form>
</div>

<style>
	.new-profile-page {
		max-width: 600px;
		margin: 0 auto;
	}
	
	.page-header {
		margin-bottom: var(--space-2xl);
	}
	
	.back-link {
		display: inline-flex;
		align-items: center;
		gap: var(--space-sm);
		color: var(--color-text-secondary);
		margin-bottom: var(--space-lg);
	}
	
	.back-link:hover {
		color: var(--color-primary);
	}
	
	.page-header h1 {
		font-size: 2rem;
		margin-bottom: var(--space-xs);
	}
	
	.page-header p {
		color: var(--color-text-secondary);
	}
	
	.profile-form {
		padding: var(--space-2xl);
	}
	
	.form-group {
		margin-bottom: var(--space-xl);
	}
	
	.form-group label {
		display: block;
		font-size: 0.875rem;
		font-weight: 500;
		margin-bottom: var(--space-sm);
		color: var(--color-text-secondary);
	}
	
	.form-group input,
	.form-group textarea {
		width: 100%;
		padding: var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text);
		font-size: 1rem;
	}
	
	.form-group input:focus,
	.form-group textarea:focus {
		outline: none;
		border-color: var(--color-primary);
	}
	
	.slug-input {
		display: flex;
		align-items: center;
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		overflow: hidden;
	}
	
	.slug-prefix {
		padding: var(--space-md);
		background: var(--color-bg-secondary);
		color: var(--color-text-muted);
		font-size: 0.875rem;
		white-space: nowrap;
	}
	
	.slug-input input {
		border: none;
		border-radius: 0;
	}
	
	.input-hint {
		display: block;
		font-size: 0.75rem;
		color: var(--color-text-muted);
		margin-top: var(--space-xs);
	}
	
	.form-actions {
		display: flex;
		gap: var(--space-md);
		justify-content: flex-end;
		margin-top: var(--space-2xl);
	}
	
	.alert {
		padding: var(--space-md);
		border-radius: var(--radius-md);
		margin-bottom: var(--space-lg);
		display: flex;
		align-items: center;
		gap: var(--space-sm);
	}
	
	.alert-error {
		background: rgba(239, 68, 68, 0.1);
		border: 1px solid rgba(239, 68, 68, 0.3);
		color: #ef4444;
	}
	
	.spin {
		animation: spin 1s linear infinite;
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
