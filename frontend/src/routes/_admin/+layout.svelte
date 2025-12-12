<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth } from '$lib/api';
	import { onMount } from 'svelte';
	
	let { children } = $props();
	
	let user = $state<{ username: string; email: string; is_admin: boolean } | null>(null);
	let sidebarOpen = $state(true);
	
	onMount(async () => {
		if (!auth.isAuthenticated()) {
			goto('/login');
			return;
		}
		
		const res = await auth.getCurrentUser();
		if (res.data) {
			user = res.data as any;
			// Check if admin
			if (!user.is_admin) {
				goto('/dashboard');
			}
		}
	});
	
	async function handleLogout() {
		await auth.logout();
		goto('/');
	}
	
	const navItems = [
		{ href: '/_admin', icon: 'bi-speedometer2', label: 'Dashboard' },
		{ href: '/_admin/users', icon: 'bi-people', label: 'Users' },
		{ href: '/_admin/profiles', icon: 'bi-layers', label: 'Profiles' },
		{ href: '/dashboard', icon: 'bi-house', label: 'Back to App' }
	];
</script>

<div class="admin-layout" class:sidebar-collapsed={!sidebarOpen}>
	<!-- Sidebar -->
	<aside class="sidebar">
		<div class="sidebar-header">
			<a href="/_admin" class="logo">
				<i class="bi bi-shield-check"></i>
				<span>Admin Panel</span>
			</a>
			<button class="sidebar-toggle" onclick={() => sidebarOpen = !sidebarOpen}>
				<i class="bi {sidebarOpen ? 'bi-chevron-left' : 'bi-list'}"></i>
			</button>
		</div>
		
		<nav class="sidebar-nav">
			{#each navItems as item}
				<a 
					href={item.href} 
					class="nav-item"
					class:active={$page.url.pathname === item.href}
				>
					<i class="bi {item.icon}"></i>
					<span>{item.label}</span>
				</a>
			{/each}
		</nav>
		
		<div class="sidebar-footer">
			{#if user}
				<div class="user-info">
					<div class="avatar admin">
						<i class="bi bi-shield-lock"></i>
					</div>
					<div class="user-details">
						<span class="username">{user.username}</span>
						<span class="badge">Admin</span>
					</div>
				</div>
			{/if}
			<button class="btn btn-ghost btn-sm" onclick={handleLogout}>
				<i class="bi bi-box-arrow-left"></i>
				<span>Logout</span>
			</button>
		</div>
	</aside>
	
	<!-- Main Content -->
	<main class="main-content">
		{@render children()}
	</main>
</div>

<style>
	.admin-layout {
		display: flex;
		min-height: 100vh;
		background: var(--color-bg);
	}
	
	/* Sidebar */
	.sidebar {
		width: 260px;
		background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
		border-right: 1px solid rgba(255,255,255,0.1);
		display: flex;
		flex-direction: column;
		transition: width var(--transition-normal);
	}
	
	.sidebar-collapsed .sidebar {
		width: 72px;
	}
	
	.sidebar-collapsed .sidebar span {
		display: none;
	}
	
	.sidebar-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: var(--space-lg);
		border-bottom: 1px solid rgba(255,255,255,0.1);
	}
	
	.logo {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--color-text);
	}
	
	.logo i {
		color: #fbbf24;
		font-size: 1.5rem;
	}
	
	.sidebar-toggle {
		padding: var(--space-sm);
		background: rgba(255,255,255,0.1);
		border: 1px solid rgba(255,255,255,0.1);
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.sidebar-toggle:hover {
		background: #fbbf24;
		color: #1a1a2e;
	}
	
	/* Navigation */
	.sidebar-nav {
		flex: 1;
		padding: var(--space-md);
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
	}
	
	.nav-item {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		padding: var(--space-md) var(--space-lg);
		color: var(--color-text-secondary);
		border-radius: var(--radius-lg);
		transition: all var(--transition-fast);
	}
	
	.nav-item:hover {
		background: rgba(255,255,255,0.1);
		color: var(--color-text);
	}
	
	.nav-item.active {
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		color: #1a1a2e;
	}
	
	.nav-item i {
		font-size: 1.25rem;
	}
	
	/* Footer */
	.sidebar-footer {
		padding: var(--space-lg);
		border-top: 1px solid rgba(255,255,255,0.1);
	}
	
	.user-info {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		margin-bottom: var(--space-md);
	}
	
	.avatar {
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255,255,255,0.1);
		border-radius: 50%;
		color: var(--color-text-secondary);
		font-size: 1.5rem;
	}
	
	.avatar.admin {
		background: linear-gradient(135deg, #fbbf24, #f59e0b);
		color: #1a1a2e;
	}
	
	.user-details {
		display: flex;
		flex-direction: column;
	}
	
	.username {
		font-weight: 500;
		font-size: 0.875rem;
	}
	
	.badge {
		font-size: 0.675rem;
		color: #fbbf24;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	
	.btn-sm {
		padding: var(--space-sm) var(--space-md);
		font-size: 0.875rem;
		width: 100%;
		justify-content: center;
	}
	
	/* Main Content */
	.main-content {
		flex: 1;
		padding: var(--space-xl);
		overflow-y: auto;
	}
</style>
