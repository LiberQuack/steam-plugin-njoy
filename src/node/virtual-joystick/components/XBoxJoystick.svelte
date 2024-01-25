<script lang="ts">
    import Analog from "./Analog.svelte";
    import Dpad from "./Dpad.svelte";
    import FaceButtons from "./FaceButtons.svelte";
    import Triggers from "./Triggers.svelte";
    import Menu from "./Menu.svelte";

    let analogLX = 0
    let analogLY = 0
    let analogRX = 0
    let analogRY = 0

    let onAnalogMoveL = (val: {x:number, y:number}) => {
      analogLX = val.x;
      analogLY = val.y;
    }
    let onAnalogMoveR = (val: {x:number, y:number}) => {
      analogRX = val.x;
      analogRY = val.y;
    }

</script>

<div style="display: flex; gap: 1em; padding: 1em;">
    <Triggers/>
    <Triggers/>
</div>

<div style="display: flex; gap: 1em; padding: 1em;">
    <Dpad/>
    <FaceButtons/>
</div>

<div style="display: flex; gap: 1em; padding: 1em;">
    <div>
        <Analog onChange={onAnalogMoveL}/>
        <div style="text-align: center">L: {analogLX.toFixed(2)} {analogLY.toFixed(2)}</div>
    </div>

    <div>
        <Analog onChange={onAnalogMoveR}/>
        <div style="text-align: center">R: {analogRX.toFixed(2)} {analogRY.toFixed(2)}</div>
    </div>
</div>

<div style="display: flex; gap: 1em; padding: 1em;">
    <Menu/>
</div>