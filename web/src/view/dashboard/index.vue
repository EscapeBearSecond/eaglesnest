<template>
  <div class="view-content ">
      <div class="chat-content bg">
          <div class="serious-char bg">1</div>
          <div class="height-char bg">2</div>
          <div class="medium-char bg">3</div>
          <div class="medium-char bg">4</div>
      </div>
      <div class="task-content bg-t">
         <span>1</span>
         <span>2</span>
         <span>3</span>
         <span>4</span>
         <span>5</span>
         <span>6</span>
      </div>
      <div class="top-content bg-f">
        <div>1</div>
        <div>2</div>
      </div>
  </div>
</template>

<script setup>

</script>

<style lang="scss" scoped>
.view-content {
  display: grid;
  height: ceil(100vh - 80);
  grid-template-rows: 3fr 1fr 3fr;
  row-gap: 15px;

  .chat-content {
    height: 35vh;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr 1fr;
    column-span: 10px;
    justify-items: center;
    align-items: center;
  }

  .task-content {
    display: grid;
    height: 10vh;
    grid-template-columns: repeat(6, 1fr);
    justify-items: center;
    align-items: center;
  }

  .top-content {
    display: grid;
    height: 35vh;
    grid-template-columns: 1fr 1fr;
    justify-items: center;
    align-items: center;
  }
}

.bg {
  background-color:antiquewhite;
}

.bg-t {
  background-color: aqua;
}

.bg-f {
  background-color: aquamarine;
}
</style>
