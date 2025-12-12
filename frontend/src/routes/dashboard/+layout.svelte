<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/api';
	import { onMount } from 'svelte';
	
	let { children } = $props();
	
	let user = $state<{ username: string; email: string } | null>(null);
	let sidebarOpen = $state(true);
	
	onMount(async () => {
		if (!auth.isAuthenticated()) {
			goto('/login');
			return;
		}
		
		const res = await auth.getCurrentUser();
		if (res.data) {
			user = res.data;
		}
	});
	
	async function handleLogout() {
		await auth.logout();
		goto('/');
	}
	
	const navItems = [
		{ href: '/dashboard', icon: 'bi-house', label: 'Dashboard' },
		{ href: '/dashboard/profiles', icon: 'bi-layers', label: 'Profiles' },
		{ href: '/dashboard/links', icon: 'bi-link-45deg', label: 'Links' },
		{ href: '/dashboard/categories', icon: 'bi-folder', label: 'Categories' },
		{ href: '/dashboard/appearance', icon: 'bi-palette', label: 'Appearance' },
		{ href: '/dashboard/analytics', icon: 'bi-graph-up', label: 'Analytics' }
	];
	
	let darkMode = $state(true);
	let showAccountMenu = $state(false);
	
	onMount(() => {
		const savedTheme = localStorage.getItem('theme');
		if (savedTheme === 'light') {
			darkMode = false;
			document.documentElement.setAttribute('data-theme', 'light');
		}
	});
	
	function toggleTheme() {
		darkMode = !darkMode;
		document.documentElement.setAttribute('data-theme', darkMode ? 'dark' : 'light');
		localStorage.setItem('theme', darkMode ? 'dark' : 'light');
	}
	
	function closeMenu() {
		showAccountMenu = false;
	}
</script>

<div class="dashboard-layout" class:sidebar-collapsed={!sidebarOpen}>
	<!-- Sidebar -->
	<aside class="sidebar">
		<div class="sidebar-header">
			<a href="/" class="logo">
				<i class="bi bi-diagram-3-fill"></i>
				<span class="logo-text">LinkMy</span>
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
					title={item.label}
				>
					<i class="bi {item.icon}"></i>
					<span class="nav-label">{item.label}</span>
				</a>
			{/each}
		</nav>
		
		<div class="sidebar-footer">
			{#if user}
				<!-- Clickable User Account -->
				<button class="user-account-btn" onclick={() => showAccountMenu = !showAccountMenu}>
					<div class="avatar">
						{#if user.avatar}
							<img src="/uploads/{user.avatar}" alt={user.username} />
						{:else}
							<i class="bi bi-person-circle"></i>
						{/if}
					</div>
					<div class="user-details">
						<span class="username">{user.username}</span>
						<span class="email">{user.email}</span>
					</div>
					<i class="bi bi-chevron-up menu-arrow" class:open={showAccountMenu}></i>
				</button>
				
				<!-- Account Menu Popup -->
				{#if showAccountMenu}
					<div class="account-menu">
						<div class="menu-header">
							<span>Account Settings</span>
						</div>
						
						<!-- Theme Toggle -->
						<button class="menu-item" onclick={toggleTheme}>
							<i class="bi {darkMode ? 'bi-sun' : 'bi-moon-stars'}"></i>
							<span>{darkMode ? 'Light Mode' : 'Dark Mode'}</span>
							<span class="toggle-indicator">{darkMode ? '🌙' : '☀️'}</span>
						</button>
						
						<!-- Settings Link -->
						<a href="/dashboard/settings" class="menu-item" onclick={closeMenu}>
							<i class="bi bi-gear"></i>
							<span>Account Settings</span>
						</a>
						
						<div class="menu-divider"></div>
						
						<!-- Logout -->
						<button class="menu-item logout" onclick={handleLogout}>
							<i class="bi bi-box-arrow-left"></i>
							<span>Logout</span>
						</button>
					</div>
				{/if}
			{/if}
		</div>
	</aside>
	
	<!-- Click outside to close menu -->
	{#if showAccountMenu}
		<div class="menu-backdrop" onclick={closeMenu}></div>
	{/if}
	
	<!-- Main Content -->
	<main class="main-content">
		{@render children()}
	</main>
</div>

<style>
	.dashboard-layout {
		display: flex;
		min-height: 100vh;
	}
	
	/* Sidebar */
	.sidebar {
		width: 260px;
		background: var(--color-bg-secondary);
		border-right: 1px solid var(--color-border);
		display: flex;
		flex-direction: column;
		transition: width var(--transition-normal);
	}
	
	.sidebar-collapsed .sidebar {
		width: 72px;
	}
	
	.sidebar-collapsed .logo-text,
	.sidebar-collapsed .nav-label,
	.sidebar-collapsed .footer-label,
	.sidebar-collapsed .btn-label,
	.sidebar-collapsed .user-details {
		display: none;
	}
	
	.sidebar-collapsed .nav-item,
	.sidebar-collapsed .footer-item {
		justify-content: center;
		padding: var(--space-md);
	}
	
	.sidebar-collapsed .user-info {
		justify-content: center;
	}
	
	.sidebar-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: var(--space-lg);
		border-bottom: 1px solid var(--color-border);
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
		color: var(--color-primary);
		font-size: 1.5rem;
	}
	
	.sidebar-toggle {
		padding: var(--space-sm);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.sidebar-toggle:hover {
		background: var(--color-primary);
		color: white;
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
		background: var(--color-bg-tertiary);
		color: var(--color-text);
	}
	
	.nav-item.active {
		background: var(--gradient-primary);
		color: white;
	}
	
	.nav-item i {
		font-size: 1.25rem;
		min-width: 20px;
		text-align: center;
	}
	
	/* Footer */
	.sidebar-footer {
		padding: var(--space-md);
		border-top: 1px solid var(--color-border);
	}
	
	.footer-item {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		padding: var(--space-sm) var(--space-md);
		margin-bottom: var(--space-xs);
		color: var(--color-text-secondary);
		border-radius: var(--radius-md);
		transition: all var(--transition-fast);
		width: 100%;
		background: none;
		border: none;
		cursor: pointer;
		font-size: 0.875rem;
		text-decoration: none;
	}
	
	.footer-item:hover {
		background: var(--color-bg-tertiary);
		color: var(--color-text);
	}
	
	.footer-item.active {
		background: var(--color-bg-tertiary);
		color: var(--color-primary);
	}
	
	.footer-item i {
		font-size: 1.125rem;
		min-width: 20px;
		text-align: center;
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
		background: var(--color-bg-tertiary);
		border-radius: 50%;
		color: var(--color-text-secondary);
		font-size: 1.5rem;
	}
	
	.user-details {
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}
	
	.username {
		font-weight: 500;
		font-size: 0.875rem;
	}
	
	.email {
		font-size: 0.75rem;
		color: var(--color-text-muted);
		text-overflow: ellipsis;
		overflow: hidden;
		white-space: nowrap;
	}
	
	/* User Account Button */
	.user-account-btn {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		padding: var(--space-sm);
		width: 100%;
		background: none;
		border: 1px solid transparent;
		border-radius: var(--radius-lg);
		color: var(--color-text);
		cursor: pointer;
		transition: all var(--transition-fast);
		text-align: left;
	}
	
	.user-account-btn:hover {
		background: var(--color-bg-tertiary);
		border-color: var(--color-border);
	}
	
	.user-account-btn .avatar {
		flex-shrink: 0;
	}
	
	.user-account-btn .avatar img {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		object-fit: cover;
	}
	
	.menu-arrow {
		margin-left: auto;
		font-size: 0.875rem;
		color: var(--color-text-muted);
		transition: transform var(--transition-fast);
	}
	
	.menu-arrow.open {
		transform: rotate(180deg);
	}
	
	/* Account Menu Popup */
	.account-menu {
		position: absolute;
		bottom: 100%;
		left: var(--space-md);
		right: var(--space-md);
		margin-bottom: var(--space-sm);
		background: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
		overflow: hidden;
		z-index: 200;
		animation: slideUp 0.2s ease-out;
	}
	
	@keyframes slideUp {
		from { opacity: 0; transform: translateY(10px); }
		to { opacity: 1; transform: translateY(0); }
	}
	
	.menu-header {
		padding: var(--space-md) var(--space-lg);
		border-bottom: 1px solid var(--color-border);
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	
	.menu-item {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		padding: var(--space-md) var(--space-lg);
		width: 100%;
		background: none;
		border: none;
		color: var(--color-text);
		font-size: 0.875rem;
		cursor: pointer;
		transition: background var(--transition-fast);
		text-decoration: none;
	}
	
	.menu-item:hover {
		background: var(--color-bg-tertiary);
	}
	
	.menu-item i {
		font-size: 1.125rem;
		color: var(--color-text-secondary);
	}
	
	.menu-item .toggle-indicator {
		margin-left: auto;
		font-size: 1rem;
	}
	
	.menu-item.logout {
		color: #ef4444;
	}
	
	.menu-item.logout i {
		color: #ef4444;
	}
	
	.menu-divider {
		height: 1px;
		background: var(--color-border);
		margin: var(--space-xs) 0;
	}
	
	/* Backdrop */
	.menu-backdrop {
		position: fixed;
		inset: 0;
		z-index: 150;
	}
	
	.sidebar-footer {
		position: relative;
	}
	
	/* Main Content */
	.main-content {
		flex: 1;
		padding: var(--space-xl);
		overflow-y: auto;
		background: var(--color-bg);
	}
	
	/* Responsive */
	@media (max-width: 768px) {
		.sidebar {
			position: fixed;
			left: 0;
			top: 0;
			bottom: 0;
			z-index: 100;
			transform: translateX(-100%);
		}
		
		.sidebar-collapsed .sidebar {
			transform: translateX(0);
			width: 260px;
		}
		
		.main-content {
			padding: var(--space-md);
		}
	}
</style>
