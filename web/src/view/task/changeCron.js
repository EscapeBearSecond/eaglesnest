import { taskForm } from './cronTask.vue';

export const changeCron = (val) => {
if (typeof val !== "string") return false;
taskForm.value.planConfig = val;
};
