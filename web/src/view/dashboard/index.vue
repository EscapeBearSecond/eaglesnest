<template>
  <div class="view-content ">
      <div class="chat-content">
          <div class="serious-char">
              <div>
                <Chart :height="height" :option="criticalOption" />
              </div>
              <div>
                <Chart :height="height" :option="highOption" />
              </div>
          </div>
          <div class="height-char">
            <div>
                <Chart :height="height" :option="mediumOption" />
              </div>
              <div>
                <Chart :height="height" :option="lowOption" />
              </div>
          </div>
      </div>
      <div class="task-content">
         <span> <el-statistic title="正在扫描" :value="taskInfo.running" /></span>
         <span>
          <span> <el-statistic title="等待扫描" :value="taskInfo.wait" /></span>
          </span>
         <span>
          <el-statistic title="扫描次数" :value="taskInfo.total" />
         </span>
         <span>
          <el-statistic title="漏洞总数" :value="vulnInfo.total">
            <!-- <template #suffix>
              <el-icon style="vertical-align: -0.125em">
                <ChatLineRound />
              </el-icon>
            </template> -->
          </el-statistic>
         </span>
         <span>
          <el-statistic title="目标总数" :value="taskInfo.targetNum" />
        </span>
        <span>
          <el-statistic title="漏洞类型" :value="vulnInfo.kindNum" />
        </span>
      </div>
      <div class="top-content">
        <div>
          <gva-card title="高危资产" custom-class="col-span-1 md:col-span-3 row-span-2">
            <gva-table />
          </gva-card>
        </div>
        <div>
          <gva-card title="漏洞排行" custom-class="col-span-1 md:col-span-3 row-span-2">
            <gva-plugin-table />
          </gva-card>
        </div>
      </div>
  </div>
</template>

<script setup>
import { taaskStatistics, vulnStatistics } from '@/api/index.js'
import { GvaPluginTable, GvaTable, GvaChart, GvaWiki , GvaNotice , GvaQuickLink , GvaCard , GvaBanner } from "./componenst"
import { ref, reactive } from 'vue'
import { useTransition } from '@vueuse/core'
import Chart from "@/components/charts/index.vue";
defineOptions({
  name: 'Index'
})


const source = ref(0)
const outputValue = useTransition(source, {
  duration: 1500,
})

const taskInfo = ref({
  failed: 0,
  running: 0,
  stopped: 0,
  success: 0,
  targetNum: 0,
  total:0
})

const vulnInfo = ref({
  critical:0,
  high:0,
  kindNum:0,
  low:0,
  medium:0,
  total:0,
})

const criticalOption =  ref({
    tooltip: {
      trigger: 'item'
    },
    color: [
      '#F81505'
    ],
    legend: {
      top: '5%',
      left: 'center',
      textStyle: {
        color: '#F81505',
      },
    },
    series: [
      {
        name: '超危',
        type: 'pie',
        radius: ['20%', '50%'],
        center: ['50%', '50%'],
        startAngle: 180,
        endAngle: 360,
        label: {
              show: true,
              position: 'outside',
              formatter: '{b}: {c}',
              fontSize: 12,
              color: '#333',
        },
        data: [
          { value: 1048, name: '超危', itemStyle: { color: '#F81505'}},
        ]
      }
    ]
})

const highOption =  ref({
    tooltip: {
      trigger: 'item'
    },
    color: [
      '#FC6C05'
    ],
    legend: {
      top: '5%',
      left: 'center',
      textStyle: {
        color: '#FC6C05',
      },
    },
    series: [
      {
        name: '高危',
        type: 'pie',
        radius: ['20%', '50%'],
        center: ['50%', '50%'],
        // startAngle: 180,
        // endAngle: 360,
        label: {
              show: true,
              position: 'outside',
              formatter: '{b}: {c}',
              fontSize: 12,
              color: '#333',
        },
        data: [
          { value: 1048, name: '高危', itemStyle: { color: '#FC6C05'}},
        ]
      }
    ]
})

const mediumOption =  ref({
    tooltip: {
      trigger: 'item'
    },
    color: [
      '#FC9D04'
    ],
    legend: {
      top: '5%',
      left: 'center',
      textStyle: {
        color: '#FC9D04',
      },
    },
    series: [
      {
        name: '中危',
        type: 'pie',
        radius: ['20%', '50%'],
        center: ['50%', '50%'],
        // startAngle: 180,
        // endAngle: 360,
        label: {
              show: true,
              position: 'outside',
              formatter: '{b}: {c}',
              fontSize: 12,
              color: '#333',
        },
        data: [
          { value: 1048, name: '中危', itemStyle: { color: '#FC9D04'}},
        ]
      }
    ]
})

const lowOption =  ref({
    tooltip: {
      trigger: 'item'
    },
    color: [
      '#A9CB21'
    ],
    legend: {
      top: '5%',
      left: 'center',
      textStyle: {
        color: '#A9CB21',
      },
    },
    series: [
      {
        name: '低危',
        type: 'pie',
        radius: ['20%', '50%'],
        center: ['50%', '50%'],
        // startAngle: 180,
        // endAngle: 360,
        label: {
              show: true,
              position: 'outside',
              formatter: '{b}: {c}',
              fontSize: 12,
              color: '#333',
        },
        data: [
          { value: 1048, name: '低危', itemStyle: { color: '#A9CB21'}},
        ]
      }
    ]
})

const initPage = async () => {
  const taskData = await taaskStatistics({})
  taskInfo.value = taskData.data
  const vulnData = await vulnStatistics({})
  vulnInfo.value = vulnData.data
  criticalOption.value.series[0].data[0].value = vulnData.data.critical
  highOption.value.series[0].data[0].value = vulnData.data.high
  mediumOption.value.series[0].data[0].value = vulnData.data.medium
  lowOption.value.series[0].data[0].value = vulnData.data.low
}

initPage()

const height = ref('328px')

source.value = 172000
</script>

<style lang="scss" scoped>
.view-content {
  display: grid;
  height: ceil(100vh - 80);
  grid-template-rows: 3fr 1fr 3fr;
  row-gap: 15px;

  .chat-content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    column-span: 10px;
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
    align-items: center;
  }
}

.serious-char, .height-char {
  display: grid;
  grid-template-columns: 1fr 1fr;
}
</style>
