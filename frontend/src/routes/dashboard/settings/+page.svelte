<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/api';
	import { onMount } from 'svelte';
	
	let user = $state<{ username: string; email: string } | null>(null);
	let loading = $state(true);
	let saving = $state(false);
	let message = $state('');
	
	// Form states
	let username = $state('');
	let email = $state('');
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	
	onMount(async () => {
		const res = await auth.getCurrentUser();
		if (res.data) {
			user = res.data;
			username = res.data.username;
			email = res.data.email;
		}
		loading = false;
	});
	
	async function updateProfile() {
		saving = true;
		message = '';
		
		const token = localStorage.getItem('access_token');
		const res = await fetch('http://localhost:3000/api/v1/me', {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify({ username, email })
		});
		
		if (res.ok) {
			message = 'Profile updated!';
		} else {
			message = 'Failed to update profile';
		}
		
		saving = false;
	}
	
	async function handleLogout() {
		await auth.logout();
		goto('/');
	}
	
	async function deleteAccount() {
		if (!confirm('Are you sure you want to delete your account? This cannot be undone.')) return;
		if (!confirm('Really delete? All your data will be lost forever.')) return;
		
		// TODO: Implement delete account API
		alert('Account deletion not yet implemented');
	}
</script>

<svelte:head>
	<title>Settings - LinkMy</title>
</svelte:head>

<div class="settings-page">
	<header class="page-header">
		<h1>Settings</h1>
		<p>Manage your account</p>
	</header>
	
	{#if loading}
		<div class="loading-state">
			<div class="loader"></div>
		</div>
	{:else}
		{#if message}
			<div class="message-toast">{message}</div>
		{/if}
		
		<!-- Profile Section -->
		<section class="settings-section card">
			<h2><i class="bi bi-person"></i> Profile</h2>
			
			<form onsubmit={(e) => { e.preventDefault(); updateProfile(); }}>
				<div class="form-group">
					<label for="username">Username</label>
					<input type="text" id="username" bind:value={username} />
				</div>
				
				<div class="form-group">
					<label for="email">Email</label>
					<input type="email" id="email" bind:value={email} />
				</div>
				
				<button type="submit" class="btn btn-primary" disabled={saving}>
					{#if saving}
						<i class="bi bi-arrow-repeat spin"></i>
						Saving...
					{:else}
						<i class="bi bi-check-lg"></i>
						Save Changes
					{/if}
				</button>
			</form>
		</section>
		
		<!-- Password Section -->
		<section class="settings-section card">
			<h2><i class="bi bi-lock"></i> Change Password</h2>
			
			<form onsubmit={(e) => { e.preventDefault(); }}>
				<div class="form-group">
					<label for="currentPassword">Current Password</label>
					<input type="password" id="currentPassword" bind:value={currentPassword} />
				</div>
				
				<div class="form-group">
					<label for="newPassword">New Password</label>
					<input type="password" id="newPassword" bind:value={newPassword} />
				</div>
				
				<div class="form-group">
					<label for="confirmPassword">Confirm New Password</label>
					<input type="password" id="confirmPassword" bind:value={confirmPassword} />
				</div>
				
				<button type="submit" class="btn btn-secondary">
					<i class="bi bi-key"></i>
					Update Password
				</button>
			</form>
		</section>
		
		<!-- Session Section -->
		<section class="settings-section card">
			<h2><i class="bi bi-box-arrow-left"></i> Session</h2>
			
			<p class="section-desc">Sign out of your account on this device.</p>
			
			<button class="btn btn-secondary" onclick={handleLogout}>
				<i class="bi bi-box-arrow-left"></i>
				Sign Out
			</button>
		</section>
		
		<!-- Danger Zone -->
		<section class="settings-section card danger-zone">
			<h2><i class="bi bi-exclamation-triangle"></i> Danger Zone</h2>
			
			<p class="section-desc">Once you delete your account, there is no going back. All your data will be permanently deleted.</p>
			
			<button class="btn btn-danger" onclick={deleteAccount}>
				<i class="bi bi-trash"></i>
				Delete Account
			</button>
		</section>
	{/if}
</div>

<style>
	.settings-page {
		max-width: 700px;
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
	
	.message-toast {
		background: var(--color-success);
		color: white;
		padding: var(--space-md) var(--space-xl);
		border-radius: var(--radius-lg);
		margin-bottom: var(--space-xl);
		animation: slideUp 0.3s ease;
	}
	
	.settings-section {
		padding: var(--space-xl);
		margin-bottom: var(--space-lg);
	}
	
	.settings-section h2 {
		font-size: 1rem;
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		margin-bottom: var(--space-xl);
		color: var(--color-text-secondary);
	}
	
	.settings-section h2 i {
		color: var(--color-primary);
	}
	
	.section-desc {
		color: var(--color-text-secondary);
		margin-bottom: var(--space-lg);
		font-size: 0.875rem;
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
	
	.form-group input {
		width: 100%;
		padding: var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text);
		font-size: 1rem;
	}
	
	.form-group input:focus {
		outline: none;
		border-color: var(--color-primary);
	}
	
	.danger-zone {
		border-color: rgba(239, 68, 68, 0.3);
	}
	
	.danger-zone h2 {
		color: #ef4444;
	}
	
	.danger-zone h2 i {
		color: #ef4444;
	}
	
	.btn-danger {
		background: #ef4444;
		color: white;
	}
	
	.btn-danger:hover {
		background: #dc2626;
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
	
	.spin { animation: spin 1s linear infinite; }
	
	@keyframes spin { to { transform: rotate(360deg); } }
	@keyframes slideUp { from { transform: translateY(20px); opacity: 0; } }
</style>
