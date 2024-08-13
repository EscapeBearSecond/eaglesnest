<template>
  <div class="view-content ">
      <div class="chat-content bg">
          <div class="serious-char bg">1</div>
          <div class="height-char bg">2</div>
          <div class="medium-char bg">3</div>
          <div class="medium-char bg">4</div>
      </div>
      <div class="task-content">
         <span> <el-statistic title="Daily active users" :value="268500" /></span>
         <span>
            <el-statistic :value="138">
              <template #title>
                <div style="display: inline-flex; align-items: center">
                  Ratio of men to women
                  <el-icon style="margin-left: 4px" :size="12">
                    <Male />
                  </el-icon>
                </div>
              </template>
              <template #suffix>/100</template>
            </el-statistic>
          </span>
         <span>
          <el-statistic title="Total Transactions" :value="outputValue" />
         </span>
         <span>
          <el-statistic title="Feedback number" :value="562">
            <template #suffix>
              <el-icon style="vertical-align: -0.125em">
                <ChatLineRound />
              </el-icon>
            </template>
          </el-statistic>
         </span>
         <span>
          <el-statistic title="Daily active users" :value="268500" />
        </span>
        <span>
          <el-statistic title="Daily active users" :value="268500" />
        </span>
      </div>
      <div class="top-content bg-f">
        <div>1</div>
        <div>2</div>
      </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useTransition } from '@vueuse/core'
defineOptions({
  name: 'Index'
})


const source = ref(0)
const outputValue = useTransition(source, {
  duration: 1500,
})

source.value = 172000
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
