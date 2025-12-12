<script lang="ts">
	// Simple Icon, Name, Category
	const iconList = [
		// Social
		{ icon: 'bi-instagram', name: 'Instagram', cat: 'social' },
		{ icon: 'bi-youtube', name: 'YouTube', cat: 'social' },
		{ icon: 'bi-tiktok', name: 'TikTok', cat: 'social' },
		{ icon: 'bi-facebook', name: 'Facebook', cat: 'social' },
		{ icon: 'bi-twitter-x', name: 'X (Twitter)', cat: 'social' },
		{ icon: 'bi-linkedin', name: 'LinkedIn', cat: 'social' },
		{ icon: 'bi-threads', name: 'Threads', cat: 'social' },
		{ icon: 'bi-pinterest', name: 'Pinterest', cat: 'social' },
		{ icon: 'bi-snapchat', name: 'Snapchat', cat: 'social' },
		{ icon: 'bi-discord', name: 'Discord', cat: 'social' },
		{ icon: 'bi-twitch', name: 'Twitch', cat: 'social' },
		{ icon: 'bi-telegram', name: 'Telegram', cat: 'social' },
		{ icon: 'bi-whatsapp', name: 'WhatsApp', cat: 'social' },
		{ icon: 'bi-reddit', name: 'Reddit', cat: 'social' },
		{ icon: 'bi-github', name: 'GitHub', cat: 'social' },
		{ icon: 'bi-spotify', name: 'Spotify', cat: 'social' },
		
		// Business
		{ icon: 'bi-globe', name: 'Website', cat: 'business' },
		{ icon: 'bi-envelope', name: 'Email', cat: 'business' },
		{ icon: 'bi-telephone', name: 'Phone', cat: 'business' },
		{ icon: 'bi-geo-alt', name: 'Location', cat: 'business' },
		{ icon: 'bi-calendar', name: 'Calendar', cat: 'business' },
		{ icon: 'bi-briefcase', name: 'Portfolio', cat: 'business' },
		{ icon: 'bi-file-earmark-text', name: 'Resume/CV', cat: 'business' },
		{ icon: 'bi-person-vcard', name: 'Contact', cat: 'business' },
		
		// Shopping
		{ icon: 'bi-shop', name: 'Shop', cat: 'shopping' },
		{ icon: 'bi-cart', name: 'Cart', cat: 'shopping' },
		{ icon: 'bi-bag', name: 'Bag', cat: 'shopping' },
		{ icon: 'bi-gift', name: 'Gift', cat: 'shopping' },
		{ icon: 'bi-credit-card', name: 'Payment', cat: 'shopping' },
		{ icon: 'bi-cash', name: 'Donate', cat: 'shopping' },
		{ icon: 'bi-cup-hot', name: 'Coffee', cat: 'shopping' },
		
		// Media
		{ icon: 'bi-play-circle', name: 'Video', cat: 'media' },
		{ icon: 'bi-music-note', name: 'Music', cat: 'media' },
		{ icon: 'bi-mic', name: 'Podcast', cat: 'media' },
		{ icon: 'bi-camera', name: 'Photo', cat: 'media' },
		{ icon: 'bi-book', name: 'Blog', cat: 'media' },
		{ icon: 'bi-journal', name: 'Article', cat: 'media' },
		{ icon: 'bi-newspaper', name: 'News', cat: 'media' },
		
		// General
		{ icon: 'bi-link-45deg', name: 'Link', cat: 'general' },
		{ icon: 'bi-star', name: 'Featured', cat: 'general' },
		{ icon: 'bi-heart', name: 'Favorite', cat: 'general' },
		{ icon: 'bi-lightning', name: 'Flash', cat: 'general' },
		{ icon: 'bi-fire', name: 'Hot', cat: 'general' },
		{ icon: 'bi-rocket', name: 'Launch', cat: 'general' },
		{ icon: 'bi-trophy', name: 'Award', cat: 'general' },
		{ icon: 'bi-bookmark', name: 'Bookmark', cat: 'general' },
		{ icon: 'bi-chat', name: 'Chat', cat: 'general' },
		{ icon: 'bi-question-circle', name: 'FAQ', cat: 'general' },
		{ icon: 'bi-info-circle', name: 'Info', cat: 'general' }
	];
	
	const categories = [
		{ id: 'all', name: 'Semua' },
		{ id: 'social', name: 'Social Media' },
		{ id: 'business', name: 'Bisnis' },
		{ id: 'shopping', name: 'Belanja' },
		{ id: 'media', name: 'Media' },
		{ id: 'general', name: 'Umum' }
	];
	
	let { value = $bindable('bi-link-45deg'), onSelect }: { value?: string; onSelect?: (icon: string) => void } = $props();
	
	let search = $state('');
	let activeCategory = $state('all');
	let isOpen = $state(false);
	
	let filteredIcons = $derived(() => {
		return iconList.filter(icon => {
			const matchSearch = icon.name.toLowerCase().includes(search.toLowerCase());
			const matchCat = activeCategory === 'all' || icon.cat === activeCategory;
			return matchSearch && matchCat;
		});
	});
	
	function selectIcon(icon: string) {
		value = icon;
		isOpen = false;
		if (onSelect) onSelect(icon);
	}
	
	function getSelectedIcon() {
		return iconList.find(i => i.icon === value);
	}
</script>

<div class="icon-picker">
	<button type="button" class="picker-trigger" onclick={() => isOpen = !isOpen}>
		<i class="bi {value}"></i>
		<span>{getSelectedIcon()?.name || 'Select Icon'}</span>
		<i class="bi bi-chevron-down"></i>
	</button>
	
	{#if isOpen}
		<div class="picker-dropdown">
			<div class="picker-header">
				<input 
					type="text" 
					placeholder="Cari icon..." 
					bind:value={search}
					class="picker-search"
				/>
			</div>
			
			<div class="picker-categories">
				{#each categories as cat}
					<button 
						type="button"
						class="cat-btn" 
						class:active={activeCategory === cat.id}
						onclick={() => activeCategory = cat.id}
					>
						{cat.name}
					</button>
				{/each}
			</div>
			
			<div class="picker-grid">
				{#each filteredIcons() as item}
					<button 
						type="button"
						class="icon-btn" 
						class:selected={value === item.icon}
						onclick={() => selectIcon(item.icon)}
						title={item.name}
					>
						<i class="bi {item.icon}"></i>
					</button>
				{/each}
			</div>
		</div>
	{/if}
</div>

<style>
	.icon-picker {
		position: relative;
	}
	
	.picker-trigger {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		width: 100%;
		padding: var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text);
		cursor: pointer;
	}
	
	.picker-trigger i:first-child {
		font-size: 1.25rem;
		color: var(--color-primary);
	}
	
	.picker-trigger span {
		flex: 1;
		text-align: left;
	}
	
	.picker-dropdown {
		position: absolute;
		top: 100%;
		left: 0;
		right: 0;
		background: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
		z-index: 100;
		margin-top: var(--space-xs);
		max-height: 400px;
		display: flex;
		flex-direction: column;
	}
	
	.picker-header {
		padding: var(--space-md);
		border-bottom: 1px solid var(--color-border);
	}
	
	.picker-search {
		width: 100%;
		padding: var(--space-sm) var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		color: var(--color-text);
	}
	
	.picker-categories {
		display: flex;
		gap: var(--space-xs);
		padding: var(--space-sm) var(--space-md);
		border-bottom: 1px solid var(--color-border);
		overflow-x: auto;
	}
	
	.cat-btn {
		padding: var(--space-xs) var(--space-sm);
		background: transparent;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		font-size: 0.75rem;
		white-space: nowrap;
		cursor: pointer;
	}
	
	.cat-btn.active {
		background: var(--color-primary);
		border-color: var(--color-primary);
		color: white;
	}
	
	.picker-grid {
		display: grid;
		grid-template-columns: repeat(8, 1fr);
		gap: var(--space-xs);
		padding: var(--space-md);
		overflow-y: auto;
		max-height: 250px;
	}
	
	.icon-btn {
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-bg-tertiary);
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.icon-btn:hover {
		background: var(--color-bg);
		border-color: var(--color-border);
		color: var(--color-text);
	}
	
	.icon-btn.selected {
		background: var(--color-primary);
		color: white;
	}
	
	.icon-btn i {
		font-size: 1.125rem;
	}
</style>
