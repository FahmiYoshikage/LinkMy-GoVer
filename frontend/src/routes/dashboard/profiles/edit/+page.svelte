<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { profiles, type Profile } from '$lib/api';
	import { onMount } from 'svelte';
	
	let profile = $state<Profile | null>(null);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');
	let success = $state('');
	
	// Edit form
	let name = $state('');
	let slug = $state('');
	let title = $state('');
	let bio = $state('');
	
	onMount(async () => {
		const id = $page.url.searchParams.get('id');
		if (!id) {
			goto('/dashboard/profiles');
			return;
		}
		
		const res = await profiles.get(parseInt(id));
		if (res.data) {
			profile = res.data;
			name = profile.name || '';
			slug = profile.slug || '';
			title = profile.title || '';
			bio = profile.bio || '';
		}
		loading = false;
	});
	
	async function handleSubmit(e: Event) {
		e.preventDefault();
		if (!profile) return;
		
		error = '';
		success = '';
		saving = true;
		
		const res = await profiles.update(profile.id, {
			name,
			slug,
			title: title || undefined,
			bio: bio || undefined
		});
		
		saving = false;
		
		if (res.error) {
			error = res.error;
		} else {
			success = 'Profile berhasil diupdate!';
			if (res.data) profile = res.data;
		}
	}
	
	function generateSlug() {
		slug = name.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/^-|-$/g, '');
	}
</script>

<svelte:head>
	<title>Edit Profile - LinkMy</title>
</svelte:head>

<div class="edit-profile-page">
	<header class="page-header">
		<a href="/dashboard/profiles" class="back-link">
			<i class="bi bi-arrow-left"></i>
			Back to Profiles
		</a>
		<h1>Edit Profile</h1>
	</header>
	
	{#if loading}
		<div class="loading">Loading...</div>
	{:else if profile}
		{#if error}
			<div class="alert alert-error">
				<i class="bi bi-exclamation-circle"></i>
				{error}
			</div>
		{/if}
		
		{#if success}
			<div class="alert alert-success">
				<i class="bi bi-check-circle"></i>
				{success}
			</div>
		{/if}
		
		<form class="card form" onsubmit={handleSubmit}>
			<div class="form-group">
				<label for="name">Profile Name</label>
				<input
					type="text"
					id="name"
					bind:value={name}
					oninput={generateSlug}
					placeholder="My Profile"
					required
				/>
				<span class="input-hint">Nama untuk membedakan profile di dashboard</span>
			</div>
			
			<div class="form-group">
				<label for="slug">URL Slug</label>
				<div class="slug-input">
					<span class="slug-prefix">linkmy.deepkernel.site/</span>
					<input
						type="text"
						id="slug"
						bind:value={slug}
						placeholder="username"
						pattern="[a-z0-9-]+"
						required
					/>
				</div>
				<span class="input-hint">URL unik untuk profile ini (hanya huruf kecil, angka, dan strip)</span>
			</div>
			
			<div class="form-group">
				<label for="title">Title</label>
				<input
					type="text"
					id="title"
					bind:value={title}
					placeholder="Web Developer | Designer"
				/>
				<span class="input-hint">Ditampilkan di bawah foto profile</span>
			</div>
			
			<div class="form-group">
				<label for="bio">Bio</label>
				<textarea
					id="bio"
					bind:value={bio}
					placeholder="Ceritakan tentang dirimu..."
					rows="3"
				></textarea>
			</div>
			
			<div class="form-actions">
				<a href="/dashboard/profiles" class="btn btn-secondary">Cancel</a>
				<button type="submit" class="btn btn-primary" disabled={saving}>
					{#if saving}
						<i class="bi bi-arrow-repeat spin"></i>
						Saving...
					{:else}
						<i class="bi bi-check-lg"></i>
						Save Changes
					{/if}
				</button>
			</div>
		</form>
		
		<div class="preview-section">
			<h2>Preview URL</h2>
			<a href="/{slug}" target="_blank" class="preview-link">
				<i class="bi bi-box-arrow-up-right"></i>
				linkmy.deepkernel.site/{slug}
			</a>
		</div>
	{/if}
</div>

<style>
	.edit-profile-page {
		max-width: 600px;
	}
	
	.page-header {
		margin-bottom: var(--space-2xl);
	}
	
	.back-link {
		display: inline-flex;
		align-items: center;
		gap: var(--space-xs);
		color: var(--color-text-secondary);
		margin-bottom: var(--space-md);
		font-size: 0.875rem;
	}
	
	.back-link:hover {
		color: var(--color-primary);
	}
	
	.page-header h1 {
		font-size: 2rem;
	}
	
	.loading {
		text-align: center;
		padding: var(--space-3xl);
		color: var(--color-text-muted);
	}
	
	.form {
		padding: var(--space-2xl);
		margin-bottom: var(--space-xl);
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
	
	.input-hint {
		display: block;
		font-size: 0.75rem;
		color: var(--color-text-muted);
		margin-top: var(--space-xs);
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
		font-size: 0.75rem;
		white-space: nowrap;
	}
	
	.slug-input input {
		border: none;
		border-radius: 0;
	}
	
	.form-actions {
		display: flex;
		gap: var(--space-md);
		justify-content: flex-end;
		margin-top: var(--space-xl);
	}
	
	.preview-section {
		padding: var(--space-xl);
		background: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
	}
	
	.preview-section h2 {
		font-size: 1rem;
		margin-bottom: var(--space-md);
		color: var(--color-text-secondary);
	}
	
	.preview-link {
		display: inline-flex;
		align-items: center;
		gap: var(--space-sm);
		color: var(--color-primary);
		font-size: 1.125rem;
	}
	
	/* Alerts */
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
	
	.alert-success {
		background: rgba(16, 185, 129, 0.1);
		border: 1px solid rgba(16, 185, 129, 0.3);
		color: #10b981;
	}
	
	.spin {
		animation: spin 1s linear infinite;
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
