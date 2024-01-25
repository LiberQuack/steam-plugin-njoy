<script lang="ts">
  import {onMount} from 'svelte';

  export let onChange: (value: { x: number, y: number }) => void

  let analogContainer: HTMLDivElement;
  let analogStick: HTMLDivElement;
  let isDragging = false;

  const handleStart = (e) => {
    isDragging = true;
    e.preventDefault();
    updateStickPosition(e);
  };

  const handleEnd = () => {
    isDragging = false;
    resetStickPosition();
    onChange({x: 0, y: 0});
  };

  const handleMove = (e) => {
    if (isDragging) {
      e.preventDefault();
      updateStickPosition(e);
    }
  };

  const updateStickPosition = (e) => {
    let clientX, clientY;

    if (e.targetTouches) {
      clientX = e.targetTouches[0].clientX;
      clientY = e.targetTouches[0].clientY;
    } else {
      clientX = e.clientX;
      clientY = e.clientY;
    }

    const rect = analogContainer.getBoundingClientRect();
    const x = clientX - rect.left;
    const y = clientY - rect.top;
    const centerX = rect.width / 2;
    const centerY = rect.height / 2;
    const dx = x - centerX;
    const dy = y - centerY;
    const distance = Math.sqrt(dx * dx + dy * dy);
    const maxDistance = 30;
    const distanceRatio = Math.min(distance / maxDistance, 1);

    // Normalize values to range -1 to 1
    const normalizedX = Math.max(-1, Math.min((dx * distanceRatio) / maxDistance, 1));
    const normalizedY = Math.max(-1, Math.min((dy * distanceRatio) / maxDistance * -1, 1)); // Invert Y axis

    onChange({x: normalizedX, y: normalizedY});

    analogStick.style.transform = `translate(-50%, -50%) translate(${dx * distanceRatio}px, ${dy * distanceRatio}px)`;
  };

  const resetStickPosition = () => {
    analogStick.style.transform = 'translate(-50%, -50%)';
  };

  onMount(() => {
    analogContainer.addEventListener('mousemove', handleMove);
    analogContainer.addEventListener('mouseup', handleEnd);
    analogContainer.addEventListener('touchmove', handleMove);
    analogContainer.addEventListener('touchend', handleEnd);
    analogContainer.addEventListener('touchcancel', handleEnd);

    return () => {
      analogContainer.removeEventListener('mousemove', handleMove);
      analogContainer.removeEventListener('mouseup', handleEnd);
      analogContainer.removeEventListener('touchmove', handleMove);
      analogContainer.removeEventListener('touchend', handleEnd);
      analogContainer.removeEventListener('touchcancel', handleEnd);
    };
  });
</script>

<div bind:this={analogContainer}
     class="analog-container controller-container"
     on:mousedown={handleStart}
     on:touchstart={handleStart}>
    <div bind:this={analogStick} class="analog-stick button"></div>
</div>
