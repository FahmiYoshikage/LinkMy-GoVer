<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/api';
	
	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);
	
	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;
		
		const res = await auth.login(email, password);
		
		if (res.error) {
			error = res.error;
			loading = false;
			return;
		}
		
		goto('/dashboard');
	}
</script>

<svelte:head>
	<title>Login - LinkMy</title>
</svelte:head>

<div class="auth-page">
	<div class="auth-container">
		<div class="auth-card card">
			<div class="auth-header">
				<a href="/" class="logo">
					<i class="bi bi-diagram-3-fill"></i>
					<span>LinkMy</span>
				</a>
				<h1>Welcome Back</h1>
				<p>Sign in to your account</p>
			</div>
			
			<form onsubmit={handleSubmit}>
				{#if error}
					<div class="alert alert-error">
						<i class="bi bi-exclamation-circle"></i>
						{error}
					</div>
				{/if}
				
				<div class="form-group">
					<label for="email">Email</label>
					<div class="input-wrapper">
						<i class="bi bi-envelope"></i>
						<input
							type="email"
							id="email"
							bind:value={email}
							placeholder="you@example.com"
							required
						/>
					</div>
				</div>
				
				<div class="form-group">
					<label for="password">Password</label>
					<div class="input-wrapper">
						<i class="bi bi-lock"></i>
						<input
							type="password"
							id="password"
							bind:value={password}
							placeholder="••••••••"
							required
						/>
					</div>
				</div>
				
				<button type="submit" class="btn btn-primary btn-full" disabled={loading}>
					{#if loading}
						<i class="bi bi-arrow-repeat spin"></i>
						Signing in...
					{:else}
						<i class="bi bi-box-arrow-in-right"></i>
						Sign In
					{/if}
				</button>
			</form>
			
			<div class="auth-footer">
				<p>Don't have an account? <a href="/register">Sign up</a></p>
			</div>
		</div>
	</div>
	
	<div class="auth-bg"></div>
</div>

<style>
	.auth-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-lg);
		position: relative;
	}
	
	.auth-bg {
		position: fixed;
		inset: 0;
		background: var(--gradient-primary);
		opacity: 0.1;
		z-index: -1;
	}
	
	.auth-container {
		width: 100%;
		max-width: 420px;
	}
	
	.auth-card {
		padding: var(--space-2xl);
	}
	
	.auth-header {
		text-align: center;
		margin-bottom: var(--space-2xl);
	}
	
	.logo {
		display: inline-flex;
		align-items: center;
		gap: var(--space-sm);
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--color-text);
		margin-bottom: var(--space-lg);
	}
	
	.logo i {
		color: var(--color-primary);
	}
	
	.auth-header h1 {
		font-size: 1.5rem;
		margin-bottom: var(--space-xs);
	}
	
	.auth-header p {
		color: var(--color-text-secondary);
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
	
	.input-wrapper {
		position: relative;
	}
	
	.input-wrapper i {
		position: absolute;
		left: var(--space-md);
		top: 50%;
		transform: translateY(-50%);
		color: var(--color-text-muted);
	}
	
	.input-wrapper input {
		width: 100%;
		padding: var(--space-md) var(--space-lg);
		padding-left: 2.75rem;
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text);
		font-size: 1rem;
		transition: all var(--transition-fast);
	}
	
	.input-wrapper input:focus {
		outline: none;
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.2);
	}
	
	.input-wrapper input::placeholder {
		color: var(--color-text-muted);
	}
	
	.btn-full {
		width: 100%;
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
	
	.auth-footer {
		text-align: center;
		margin-top: var(--space-xl);
		padding-top: var(--space-xl);
		border-top: 1px solid var(--color-border);
		color: var(--color-text-secondary);
	}
	
	.spin {
		animation: spin 1s linear infinite;
	}
	
	@keyframes spin {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}
</style>
