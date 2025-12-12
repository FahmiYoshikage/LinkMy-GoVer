<script lang="ts">
	let { 
		isOpen = $bindable(false), 
		slug,
		profileName 
	}: { 
		isOpen?: boolean; 
		slug: string;
		profileName: string;
	} = $props();
	
	const profileUrl = `https://linkmy.deepkernel.site/${slug}`;
	
	// Generate QR code using API
	const qrUrl = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(profileUrl)}&bgcolor=ffffff&color=000000&format=svg`;
	
	function close() {
		isOpen = false;
	}
	
	function copyLink() {
		navigator.clipboard.writeText(profileUrl);
		// Show toast or feedback
	}
	
	function shareTwitter() {
		window.open(`https://twitter.com/intent/tweet?url=${encodeURIComponent(profileUrl)}&text=${encodeURIComponent(`Check out ${profileName}'s links!`)}`, '_blank');
	}
	
	function shareFacebook() {
		window.open(`https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(profileUrl)}`, '_blank');
	}
	
	function shareWhatsApp() {
		window.open(`https://wa.me/?text=${encodeURIComponent(`${profileName}: ${profileUrl}`)}`, '_blank');
	}
	
	function shareTelegram() {
		window.open(`https://t.me/share/url?url=${encodeURIComponent(profileUrl)}&text=${encodeURIComponent(`Check out ${profileName}'s links!`)}`, '_blank');
	}
	
	function downloadQR() {
		const link = document.createElement('a');
		link.href = qrUrl.replace('svg', 'png');
		link.download = `${slug}-qr.png`;
		link.click();
	}
</script>

{#if isOpen}
	<div class="modal-overlay" onclick={close}>
		<div class="share-modal" onclick={(e) => e.stopPropagation()}>
			<button class="close-btn" onclick={close}>
				<i class="bi bi-x-lg"></i>
			</button>
			
			<div class="modal-header">
				<h2>Share Profile</h2>
			</div>
			
			<!-- QR Code -->
			<div class="qr-section">
				<div class="qr-code">
					<img src={qrUrl} alt="QR Code for {slug}" />
				</div>
				<button class="btn btn-secondary btn-sm" onclick={downloadQR}>
					<i class="bi bi-download"></i>
					Download QR
				</button>
			</div>
			
			<!-- Copy Link -->
			<div class="link-section">
				<div class="link-box">
					<span>{profileUrl}</span>
					<button class="copy-btn" onclick={copyLink}>
						<i class="bi bi-clipboard"></i>
					</button>
				</div>
			</div>
			
			<!-- Social Sharing -->
			<div class="social-section">
				<p>Share to:</p>
				<div class="social-buttons">
					<button class="social-btn twitter" onclick={shareTwitter}>
						<i class="bi bi-twitter-x"></i>
					</button>
					<button class="social-btn facebook" onclick={shareFacebook}>
						<i class="bi bi-facebook"></i>
					</button>
					<button class="social-btn whatsapp" onclick={shareWhatsApp}>
						<i class="bi bi-whatsapp"></i>
					</button>
					<button class="social-btn telegram" onclick={shareTelegram}>
						<i class="bi bi-telegram"></i>
					</button>
				</div>
			</div>
			
			<!-- Join CTA -->
			<div class="cta-section">
				<div class="cta-divider">
					<span>Mau punya link page sendiri?</span>
				</div>
				<a href="/register" class="btn btn-primary btn-full">
					<i class="bi bi-rocket-takeoff"></i>
					Buat Gratis di LinkMy
				</a>
			</div>
		</div>
	</div>
{/if}

<style>
	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.8);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		z-index: 1000;
		animation: fadeIn 0.2s ease-out;
	}
	
	.share-modal {
		position: relative;
		width: 100%;
		max-width: 380px;
		background: #1a1a2e;
		border-radius: 24px;
		padding: 2rem;
		color: white;
		animation: slideUp 0.3s ease-out;
	}
	
	.close-btn {
		position: absolute;
		top: 1rem;
		right: 1rem;
		width: 36px;
		height: 36px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255,255,255,0.1);
		border: none;
		border-radius: 50%;
		color: white;
		cursor: pointer;
		transition: background 0.2s;
	}
	
	.close-btn:hover {
		background: rgba(255,255,255,0.2);
	}
	
	.modal-header {
		text-align: center;
		margin-bottom: 1.5rem;
	}
	
	.modal-header h2 {
		font-size: 1.25rem;
		font-weight: 600;
	}
	
	.qr-section {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		margin-bottom: 1.5rem;
	}
	
	.qr-code {
		width: 180px;
		height: 180px;
		background: white;
		border-radius: 16px;
		padding: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.qr-code img {
		width: 100%;
		height: 100%;
	}
	
	.link-section {
		margin-bottom: 1.5rem;
	}
	
	.link-box {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		background: rgba(255,255,255,0.1);
		border-radius: 12px;
	}
	
	.link-box span {
		flex: 1;
		font-size: 0.875rem;
		color: rgba(255,255,255,0.8);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.copy-btn {
		width: 36px;
		height: 36px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #667eea;
		border: none;
		border-radius: 8px;
		color: white;
		cursor: pointer;
		transition: background 0.2s;
	}
	
	.copy-btn:hover {
		background: #5a6fd9;
	}
	
	.social-section {
		margin-bottom: 1.5rem;
	}
	
	.social-section p {
		font-size: 0.875rem;
		color: rgba(255,255,255,0.6);
		margin-bottom: 0.75rem;
		text-align: center;
	}
	
	.social-buttons {
		display: flex;
		justify-content: center;
		gap: 0.75rem;
	}
	
	.social-btn {
		width: 48px;
		height: 48px;
		display: flex;
		align-items: center;
		justify-content: center;
		border: none;
		border-radius: 12px;
		font-size: 1.25rem;
		color: white;
		cursor: pointer;
		transition: transform 0.2s, opacity 0.2s;
	}
	
	.social-btn:hover {
		transform: scale(1.1);
	}
	
	.social-btn.twitter { background: #1da1f2; }
	.social-btn.facebook { background: #1877f2; }
	.social-btn.whatsapp { background: #25d366; }
	.social-btn.telegram { background: #0088cc; }
	
	.cta-section {
		text-align: center;
	}
	
	.cta-divider {
		display: flex;
		align-items: center;
		gap: 1rem;
		margin-bottom: 1rem;
	}
	
	.cta-divider::before,
	.cta-divider::after {
		content: '';
		flex: 1;
		height: 1px;
		background: rgba(255,255,255,0.2);
	}
	
	.cta-divider span {
		font-size: 0.75rem;
		color: rgba(255,255,255,0.5);
		white-space: nowrap;
	}
	
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 12px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}
	
	.btn-primary {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
	}
	
	.btn-primary:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
	}
	
	.btn-secondary {
		background: rgba(255,255,255,0.1);
		color: white;
	}
	
	.btn-full {
		width: 100%;
	}
	
	.btn-sm {
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
	}
	
	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}
	
	@keyframes slideUp {
		from { transform: translateY(20px); opacity: 0; }
		to { transform: translateY(0); opacity: 1; }
	}
</style>
