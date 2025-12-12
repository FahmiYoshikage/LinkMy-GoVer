<script lang="ts">
	import { onMount } from 'svelte';
	import { profiles, themes, type Profile, type Theme } from '$lib/api';
	
	let userProfiles = $state<Profile[]>([]);
	let selectedProfileId = $state<number | null>(null);
	let theme = $state<Theme | null>(null);
	let loading = $state(true);
	let saving = $state(false);
	let saveMessage = $state('');
	
	onMount(async () => {
		const res = await profiles.getAll();
		if (res.data && res.data.length > 0) {
			userProfiles = res.data;
			selectedProfileId = res.data[0].id;
			await loadTheme();
		}
		loading = false;
	});
	
	async function loadTheme() {
		if (!selectedProfileId) return;
		const res = await themes.get(selectedProfileId);
		if (res.data) {
			theme = res.data;
		}
	}
	
	async function handleProfileChange(e: Event) {
		const select = e.target as HTMLSelectElement;
		selectedProfileId = parseInt(select.value);
		await loadTheme();
	}
	
	async function saveTheme() {
		if (!theme || !selectedProfileId) return;
		saving = true;
		
		const res = await themes.update(selectedProfileId, theme);
		if (res.data) {
			theme = res.data;
			saveMessage = 'Theme saved!';
			setTimeout(() => saveMessage = '', 3000);
		}
		
		saving = false;
	}
	
	const gradientPresets = [
		{ name: 'Purple Dream', value: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
		{ name: 'Sunset', value: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
		{ name: 'Ocean', value: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
		{ name: 'Forest', value: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)' },
		{ name: 'Fire', value: 'linear-gradient(135deg, #f12711 0%, #f5af19 100%)' },
		{ name: 'Night', value: 'linear-gradient(135deg, #232526 0%, #414345 100%)' },
		{ name: 'Cotton Candy', value: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)' },
		{ name: 'Midnight', value: 'linear-gradient(135deg, #2c3e50 0%, #3498db 100%)' },
		{ name: 'Rose', value: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)' },
		{ name: 'Peach', value: 'linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)' },
		{ name: 'Aurora', value: 'linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%)' },
		{ name: 'Neon', value: 'linear-gradient(135deg, #00f260 0%, #0575e6 100%)' }
	];
	
	const fontOptions = [
		'Inter', 'Roboto', 'Poppins', 'Open Sans', 'Lato', 'Montserrat', 'Outfit', 'Nunito',
		'Playfair Display', 'Merriweather', 'Space Grotesk', 'DM Sans'
	];
	
	const buttonStyles = [
		{ value: 'rounded', label: 'Rounded' },
		{ value: 'pill', label: 'Pill' },
		{ value: 'square', label: 'Square' },
		{ value: 'outline', label: 'Outline' },
		{ value: 'soft', label: 'Soft' }
	];
	
	const shadowOptions = [
		{ value: 'none', label: 'None' },
		{ value: 'light', label: 'Light' },
		{ value: 'medium', label: 'Medium' },
		{ value: 'heavy', label: 'Heavy' }
	];
</script>

<svelte:head>
	<title>Appearance - LinkMy</title>
</svelte:head>

<div class="appearance-page">
	<header class="page-header">
		<div class="header-left">
			<h1>Appearance</h1>
			<p>Customize your profile look</p>
		</div>
		<div class="header-actions">
			{#if userProfiles.length > 1}
				<select class="profile-select" onchange={handleProfileChange} value={selectedProfileId}>
					{#each userProfiles as profile}
						<option value={profile.id}>{profile.name}</option>
					{/each}
				</select>
			{/if}
			<button class="btn btn-primary" onclick={saveTheme} disabled={saving}>
				{#if saving}
					<i class="bi bi-arrow-repeat spin"></i>
					Saving...
				{:else}
					<i class="bi bi-check-lg"></i>
					Save Changes
				{/if}
			</button>
		</div>
	</header>
	
	{#if saveMessage}
		<div class="save-toast">{saveMessage}</div>
	{/if}
	
	{#if loading}
		<div class="loading-state">
			<div class="loader"></div>
		</div>
	{:else if theme}
		<div class="editor-grid">
			<!-- Background Section -->
			<section class="editor-section card">
				<h2><i class="bi bi-palette"></i> Background</h2>
				
				<div class="form-group">
					<label>Type</label>
					<div class="btn-group">
						<button 
							class:active={theme.bg_type === 'solid'} 
							onclick={() => theme && (theme.bg_type = 'solid')}
						>Solid</button>
						<button 
							class:active={theme.bg_type === 'gradient'} 
							onclick={() => theme && (theme.bg_type = 'gradient')}
						>Gradient</button>
						<button 
							class:active={theme.bg_type === 'image'} 
							onclick={() => theme && (theme.bg_type = 'image')}
						>Image</button>
					</div>
				</div>
				
				{#if theme.bg_type === 'solid'}
					<div class="form-group">
						<label>Color</label>
						<div class="color-input">
							<input type="color" bind:value={theme.bg_value} />
							<input type="text" bind:value={theme.bg_value} placeholder="#667eea" />
						</div>
					</div>
				{:else if theme.bg_type === 'gradient'}
					<div class="form-group">
						<label>Preset Gradients</label>
						<div class="gradient-grid">
							{#each gradientPresets as preset}
								<button
									class="gradient-btn"
									class:active={theme.bg_value === preset.value}
									style="background: {preset.value}"
									title={preset.name}
									onclick={() => theme && (theme.bg_value = preset.value)}
								></button>
							{/each}
						</div>
					</div>
				{/if}
			</section>
			
			<!-- Buttons Section -->
			<section class="editor-section card">
				<h2><i class="bi bi-ui-radios"></i> Buttons</h2>
				
				<div class="form-group">
					<label>Style</label>
					<div class="btn-group">
						{#each buttonStyles as style}
							<button 
								class:active={theme.button_style === style.value}
								onclick={() => theme && (theme.button_style = style.value)}
							>{style.label}</button>
						{/each}
					</div>
				</div>
				
				<div class="form-group">
					<label>Button Color</label>
					<div class="color-input">
						<input type="color" bind:value={theme.button_color} />
						<input type="text" bind:value={theme.button_color} placeholder="#667eea" />
					</div>
				</div>
				
				<div class="form-group">
					<label>Text Color</label>
					<div class="color-input">
						<input type="color" bind:value={theme.text_color} />
						<input type="text" bind:value={theme.text_color} placeholder="#333333" />
					</div>
				</div>
			</section>
			
			<!-- Typography Section -->
			<section class="editor-section card">
				<h2><i class="bi bi-fonts"></i> Typography</h2>
				
				<div class="form-group">
					<label>Font</label>
					<select bind:value={theme.font}>
						{#each fontOptions as font}
							<option value={font}>{font}</option>
						{/each}
					</select>
				</div>
			</section>
			
			<!-- Effects Section -->
			<section class="editor-section card">
				<h2><i class="bi bi-stars"></i> Effects</h2>
				
				<div class="toggle-group">
					<label class="toggle-label">
						<span>Animations</span>
						<input type="checkbox" bind:checked={theme.enable_animations} />
						<span class="toggle-switch"></span>
					</label>
					
					<label class="toggle-label">
						<span>Glass Effect</span>
						<input type="checkbox" bind:checked={theme.enable_glass_effect} />
						<span class="toggle-switch"></span>
					</label>
				</div>
			</section>
			
			<!-- Boxed Layout Section -->
			<section class="editor-section card full-width">
				<h2><i class="bi bi-layout-wtf"></i> Boxed Layout</h2>
				
				<label class="toggle-label">
					<span>Enable Boxed Layout</span>
					<input type="checkbox" bind:checked={theme.boxed_enabled} />
					<span class="toggle-switch"></span>
				</label>
				
				{#if theme.boxed_enabled}
					<div class="boxed-options">
						<div class="form-group">
							<label>Container Background</label>
							<div class="color-input">
								<input type="color" bind:value={theme.boxed_container_bg} />
								<input type="text" bind:value={theme.boxed_container_bg} />
							</div>
						</div>
						
						<div class="form-group">
							<label>Max Width: {theme.boxed_max_width}px</label>
							<input type="range" min="400" max="800" bind:value={theme.boxed_max_width} />
						</div>
						
						<div class="form-group">
							<label>Border Radius: {theme.boxed_radius}px</label>
							<input type="range" min="0" max="50" bind:value={theme.boxed_radius} />
						</div>
						
						<label class="toggle-label">
							<span>Container Shadow</span>
							<input type="checkbox" bind:checked={theme.boxed_shadow} />
							<span class="toggle-switch"></span>
						</label>
					</div>
				{/if}
			</section>
		</div>
		
		<!-- Preview -->
		<section class="preview-section">
			<h2>Preview</h2>
			<div class="preview-frame" style="background: {theme.bg_value || '#667eea'}">
				<div 
					class="preview-container"
					class:boxed={theme.boxed_enabled}
					style="
						{theme.boxed_enabled ? `
							background: ${theme.boxed_container_bg};
							max-width: ${theme.boxed_max_width}px;
							border-radius: ${theme.boxed_radius}px;
							${theme.boxed_shadow ? 'box-shadow: 0 20px 60px rgba(0,0,0,0.3);' : ''}
						` : ''}
						font-family: '{theme.font}', sans-serif;
						color: {theme.text_color};
					"
				>
					<div class="preview-avatar"></div>
					<h3>Preview</h3>
					<p>This is how your profile will look</p>
					
					{#each [1, 2, 3] as i}
						<button
							class="preview-btn"
							class:pill={theme.button_style === 'pill'}
							class:rounded={theme.button_style === 'rounded'}
							class:square={theme.button_style === 'square'}
							class:glass={theme.enable_glass_effect}
							style="background: {theme.enable_glass_effect ? 'rgba(255,255,255,0.15)' : theme.button_color}"
						>
							Link {i}
						</button>
					{/each}
				</div>
			</div>
		</section>
	{/if}
</div>

<style>
	.appearance-page {
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
	
	.save-toast {
		position: fixed;
		bottom: var(--space-xl);
		right: var(--space-xl);
		background: var(--color-success);
		color: white;
		padding: var(--space-md) var(--space-xl);
		border-radius: var(--radius-lg);
		animation: slideUp 0.3s ease;
		z-index: 100;
	}
	
	/* Editor Grid */
	.editor-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: var(--space-lg);
		margin-bottom: var(--space-2xl);
	}
	
	.editor-section {
		padding: var(--space-xl);
	}
	
	.editor-section.full-width {
		grid-column: 1 / -1;
	}
	
	.editor-section h2 {
		font-size: 1rem;
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		margin-bottom: var(--space-xl);
		color: var(--color-text-secondary);
	}
	
	.editor-section h2 i {
		color: var(--color-primary);
	}
	
	/* Form Elements */
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
	
	.form-group select,
	.form-group input[type="text"] {
		width: 100%;
		padding: var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		color: var(--color-text);
	}
	
	.btn-group {
		display: flex;
		gap: var(--space-xs);
	}
	
	.btn-group button {
		flex: 1;
		padding: var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		transition: all var(--transition-fast);
	}
	
	.btn-group button.active {
		background: var(--gradient-primary);
		border-color: transparent;
		color: white;
	}
	
	.color-input {
		display: flex;
		gap: var(--space-sm);
	}
	
	.color-input input[type="color"] {
		width: 48px;
		height: 48px;
		border: none;
		border-radius: var(--radius-md);
		cursor: pointer;
	}
	
	.color-input input[type="text"] {
		flex: 1;
	}
	
	.gradient-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: var(--space-sm);
	}
	
	.gradient-btn {
		aspect-ratio: 1;
		border: 2px solid transparent;
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.gradient-btn.active {
		border-color: white;
		transform: scale(1.1);
	}
	
	/* Toggle */
	.toggle-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
	}
	
	.toggle-label {
		display: flex;
		align-items: center;
		justify-content: space-between;
		cursor: pointer;
	}
	
	.toggle-label input {
		display: none;
	}
	
	.toggle-switch {
		width: 48px;
		height: 28px;
		background: var(--color-bg-tertiary);
		border-radius: 14px;
		position: relative;
		transition: background var(--transition-fast);
	}
	
	.toggle-switch::after {
		content: '';
		position: absolute;
		width: 22px;
		height: 22px;
		background: white;
		border-radius: 50%;
		top: 3px;
		left: 3px;
		transition: transform var(--transition-fast);
	}
	
	.toggle-label input:checked + .toggle-switch {
		background: var(--color-primary);
	}
	
	.toggle-label input:checked + .toggle-switch::after {
		transform: translateX(20px);
	}
	
	.boxed-options {
		margin-top: var(--space-xl);
		padding-top: var(--space-xl);
		border-top: 1px solid var(--color-border);
	}
	
	input[type="range"] {
		width: 100%;
		accent-color: var(--color-primary);
	}
	
	/* Preview */
	.preview-section {
		margin-top: var(--space-2xl);
	}
	
	.preview-section h2 {
		margin-bottom: var(--space-lg);
	}
	
	.preview-frame {
		min-height: 400px;
		border-radius: var(--radius-xl);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-xl);
	}
	
	.preview-container {
		text-align: center;
		padding: var(--space-2xl);
		width: 100%;
		max-width: 400px;
	}
	
	.preview-container.boxed {
		background: white;
	}
	
	.preview-avatar {
		width: 80px;
		height: 80px;
		background: rgba(255, 255, 255, 0.2);
		border-radius: 50%;
		margin: 0 auto var(--space-lg);
	}
	
	.preview-container h3 {
		margin-bottom: var(--space-sm);
	}
	
	.preview-container p {
		opacity: 0.7;
		margin-bottom: var(--space-xl);
	}
	
	.preview-btn {
		display: block;
		width: 100%;
		padding: var(--space-md) var(--space-lg);
		margin-bottom: var(--space-sm);
		color: white;
		border: none;
		cursor: pointer;
	}
	
	.preview-btn.pill { border-radius: 50px; }
	.preview-btn.rounded { border-radius: var(--radius-lg); }
	.preview-btn.square { border-radius: 0; }
	
	.preview-btn.glass {
		backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.2);
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
