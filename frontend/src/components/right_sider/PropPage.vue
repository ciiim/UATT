<template>
  <div class="prop-page" v-if="store.selectedAction">
    <h3 style="margin-bottom: 16px; color: black">
      属性编辑 - {{ store.selectedAction.Name }}
    </h3>

    <!-- IO 类模块 -->
    <template v-if="isIO">
      <a-form layout="vertical">
        <!-- IO 顶级字段 -->
        <a-form-item label="超时 (ms)">
          <a-input-number
            v-model:value="store.selectedAction.TypeFeatureField.TimeoutMs"
            style="width: 100%"
          />
        </a-form-item>

        <!-- IO 子模块数组 -->
        <template
          v-for="(mod, index) in store.selectedAction.TypeFeatureField.Modules"
          :key="mod.ModuleUID"
        >
          <a-divider>子模块 {{ index + 1 }}</a-divider>

          <a-form-item label="ModuleTypeID">
            <a-input-number
              v-model:value="mod.ModuleTypeID"
              disabled
              style="width: 100%"
            />
          </a-form-item>

          <!-- Fill -->
          <template v-if="mod.ModuleTypeID === 10">
            <a-form-item label="使用变量">
              <a-input v-model:value="mod.UseVar" />
            </a-form-item>
            <a-form-item label="填充长度">
              <a-input-number
                v-model:value="mod.FillLength"
                style="width: 100%"
              />
            </a-form-item>
          </template>

          <!-- Fixed -->
          <template v-else-if="mod.ModuleTypeID === 11">
            <a-form-item label="固定内容 (字节数组)">
              <a-textarea
                :value="fixedToString(mod.FixedContent)"
                @blur="(e : any)=> mod.FixedContent = fixedFromString(e.target.value)"
              />
            </a-form-item>
          </template>

          <!-- Calc -->
          <template v-else-if="mod.ModuleTypeID === 12">
            <a-form-item label="模式">
              <a-input v-model:value="mod.Mode" />
            </a-form-item>
            <a-form-item label="计算函数">
              <a-input v-model:value="mod.CalcFunc" />
            </a-form-item>
            <a-form-item label="计算时机">
              <a-input v-model:value="mod.CalcTiming" />
            </a-form-item>
            <a-form-item label="占位字节">
              <a-textarea
                :value="fixedToString(mod.PlaceholderBytes)"
                @blur="
                  (e : any) => mod.PlaceholderBytes = fixedFromString(
                    e.target.value
                  )
                "
              />
            </a-form-item>
            <!-- 多选输入模块UID -->
            <a-form-item label="输入模块 UID">
              <a-select
                v-model:value="mod.CalcInputModulesUID"
                mode="multiple"
                style="width: 100%"
                :options="calcInputOptions(mod)"
                placeholder="请选择输入模块"
              />
            </a-form-item>
          </template>

          <!-- Custom -->
          <template v-else-if="mod.ModuleTypeID === 13">
            <a-form-item label="自定义长度">
              <a-input-number
                v-model:value="mod.CustomLength"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item label="自定义内容 (字节数组)">
              <a-textarea
                :value="fixedToString(mod.CustomContent)"
                @blur="
                  (e : any)=> mod.CustomContent = fixedFromString(e.target.value)
                "
              />
            </a-form-item>
          </template>
          <a-button danger @click="removeModule(index)">删除模块</a-button>
          <Divider></Divider>
        </template>
        <div style="padding-top: 20px">
          <a-popconfirm
            title="请选择要新增的模块类型"
            @confirm="() => addModule(newModuleType)"
            @cancel="() => (newModuleType = null)"
          >
            <template #description>
              <!-- 在 Popconfirm Description 里放一个 Select -->
              <a-select v-model:value="newModuleType" style="width: 180px">
                <a-select-option :value="10">Fill</a-select-option>
                <a-select-option :value="11">Fixed</a-select-option>
                <a-select-option :value="12">Calc</a-select-option>
                <a-select-option :value="13">Custom</a-select-option>
              </a-select>
            </template>

            <a-button type="primary">新增模块</a-button>
          </a-popconfirm>
        </div>
      </a-form>
    </template>

    <!-- Control 类模块 -->
    <template v-else-if="isControl">
      <a-form layout="vertical">
        <!-- Declare -->
        <template v-if="store.selectedAction.ActionTypeID === 23">
          <a-form-item label="变量名">
            <a-input
              v-model:value="store.selectedAction.TypeFeatureField.VarName"
            />
          </a-form-item>
          <a-form-item label="变量类型">
            <a-input
              v-model:value="store.selectedAction.TypeFeatureField.VarType"
            />
          </a-form-item>
          <a-form-item label="数值(Number)">
            <a-input-number
              v-model:value="
                store.selectedAction.TypeFeatureField.VarNumberValue
              "
              style="width: 100%"
            />
          </a-form-item>
          <a-form-item label="数值(String)">
            <a-input
              v-model:value="
                store.selectedAction.TypeFeatureField.VarStringValue
              "
            />
          </a-form-item>
          <a-form-item label="字节数组值">
            <a-textarea
              :value="
                fixedToString(
                  store.selectedAction.TypeFeatureField.VarByteArrayValue
                )
              "
              @blur="
                (e : any) => {
                    if(store.selectedAction != undefined)
                        store.selectedAction.TypeFeatureField.VarByteArrayValue = fixedFromString(e.target.value)
                }
              "
            />
          </a-form-item>
        </template>

        <!-- IF -->
        <template v-else-if="store.selectedAction.ActionTypeID === 24">
          <a-form-item label="条件表达式">
            <a-input
              v-model:value="store.selectedAction.TypeFeatureField.Condition"
            />
          </a-form-item>
        </template>

        <!-- ELSE -->
        <template v-else-if="store.selectedAction.ActionTypeID === 25">
          <p>ELSE 模块无额外属性</p>
        </template>

        <!-- FOR -->
        <template v-else-if="store.selectedAction.ActionTypeID === 26">
          <a-form-item label="使用变量">
            <a-input
              v-model:value="store.selectedAction.TypeFeatureField.UseVar"
            />
          </a-form-item>
          <a-form-item label="进入条件">
            <a-input
              v-model:value="
                store.selectedAction.TypeFeatureField.EnterCondition
              "
            />
          </a-form-item>
          <a-form-item label="变量操作">
            <a-input
              v-model:value="store.selectedAction.TypeFeatureField.VarOp"
            />
          </a-form-item>
        </template>

        <!-- EndBlock -->
        <template v-else-if="store.selectedAction.ActionTypeID === 27">
          <p>EndBlock 模块无额外属性</p>
        </template>

        <!-- Label -->
        <template v-else-if="store.selectedAction.ActionTypeID === 28">
          <a-form-item label="标签名">
            <a-input
              v-model:value="store.selectedAction.TypeFeatureField.LabelName"
            />
          </a-form-item>
        </template>

        <!-- GOTO -->
        <template v-else-if="store.selectedAction.ActionTypeID === 29">
          <a-form-item label="跳转标签">
            <a-input
              v-model:value="store.selectedAction.TypeFeatureField.Label"
            />
          </a-form-item>
        </template>

        <!-- ChangeBaudrate -->
        <template v-else-if="store.selectedAction.ActionTypeID === 30">
          <a-form-item label="目标波特率">
            <a-input-number
              v-model:value="
                store.selectedAction.TypeFeatureField.TargetBaudRate
              "
              style="width: 100%"
            />
          </a-form-item>
        </template>

        <!-- Stop -->
        <template v-else-if="store.selectedAction.ActionTypeID === 31">
          <a-form-item label="停止代码">
            <a-input-number
              v-model:value="store.selectedAction.TypeFeatureField.StopCode"
              style="width: 100%"
            />
          </a-form-item>
        </template>
      </a-form>
    </template>

    <!-- Debug 类模块 -->
    <template v-else-if="isDebug">
      <a-form layout="vertical">
        <!-- Print -->
        <template v-if="store.selectedAction.ActionTypeID === 90">
          <a-form-item label="打印格式字符串">
            <a-input
              v-model:value="store.selectedAction.TypeFeatureField.PrintFmt"
            />
          </a-form-item>
        </template>

        <!-- Delay -->
        <template v-else-if="store.selectedAction.ActionTypeID === 91">
          <a-form-item label="延迟时间(ms)">
            <a-input-number
              v-model:value="store.selectedAction.TypeFeatureField.DelayMs"
              style="width: 100%"
            />
          </a-form-item>
        </template>
      </a-form>
    </template>

    <template v-else>
      <p style="color: black">未识别的 Action 类型</p>
    </template>
  </div>

  <div v-else>
    <p style="color: black">请选择一个 Action 编辑属性</p>
  </div>
</template>

<script setup lang="ts">
import { computed, watch, ref } from "vue";
import type { ConfigActionBase } from "../../types/Action";
import { useActionStore } from "../../stores/action_store";
import { parseActionTags } from "../../utils/action_utils";
import { Divider } from "ant-design-vue";

const store = useActionStore();

const newModuleType = ref<number | null>(null);

const emit = defineEmits<{
  (e: "update-tags", updatedAction: ConfigActionBase): void;
}>();

// 当 TypeFeatureField 变化时调用
watch(
  () => store.selectedAction?.TypeFeatureField,
  (newVal) => {
    if (store.selectedAction) {
      const sa = store.selectedAction as any;
      sa.tags = parseActionTags(store.selectedAction);
      store.selectedAction = sa;
    }
  },
  { deep: true }
);

const addModule = (typeId: number | null) => {
  if (!store.selectedAction || !typeId) return;

  const modules = store.selectedAction.TypeFeatureField.Modules;
  const nextUid = genNextUID(modules);

  const newModule = createModuleTemplate(typeId, nextUid);
  modules.push(newModule);

  newModuleType.value = null;
};

const genNextUID = (modules: any[]) => {
  return Date.now();
};

const createModuleTemplate = (typeId: number, uid: number) => {
  switch (typeId) {
    case 10:
      return {
        ModuleUID: uid,
        ModuleTypeID: typeId,
        FillLength: 0,
        UseVar: "",
      };
    case 11:
      return {
        ModuleUID: uid,
        ModuleTypeID: typeId,
        FixedContent: [],
      };
    case 12:
      return {
        ModuleUID: uid,
        ModuleTypeID: typeId,
        Mode: "",
        CalcFunc: "",
        CalcTiming: "",
        PlaceholderBytes: [],
        CalcInputModulesUID: [],
      };
    case 13:
      return {
        ModuleUID: uid,
        ModuleTypeID: typeId,
        CustomContent: [],
      };
    default:
      return {
        ModuleUID: uid,
        ModuleTypeID: typeId,
      };
  }
};

const isIO = computed(() => store.selectedAction?.ActionType === "IO");
const isControl = computed(
  () => store.selectedAction?.ActionType === "Control"
);
const isDebug = computed(() => store.selectedAction?.ActionType === "Debug");

// 工具函数: 将字节数组转成字符串显示
const fixedToString = (arr: number[] | undefined): string => {
  if (!arr || !Array.isArray(arr)) return "";
  return arr
    .map((n) => {
      // 转成大写十六进制，补零
      const hex = n.toString(16).toUpperCase().padStart(2, "0");
      return `0x${hex}`; // 加前缀，便于识别
    })
    .join(",");
};
const fixedFromString = (str: any) => {
  console.log(str);

  if (!str) return [];
  return str
    .split(",")
    .map((s: any) => parseInt(s))
    .filter((n: any) => !isNaN(n));
};

const removeModule = (index: number) => {
  store.selectedAction?.TypeFeatureField.Modules.splice(index, 1);
};

const moduleList = new Map([
  [10, "Fill"],
  [11, "Fixed"],
  [12, "Calc"],
  [13, "Custom"],
]);

const calcInputOptions = (currentMod: any) => {
  const modules = store.selectedAction?.TypeFeatureField?.Modules ?? [];
  return modules
    .filter((m: any) => m.ModuleUID !== currentMod.ModuleUID) // 排除自己
    .map((m: any) => ({
      label: `模块UID ${m.ModuleUID} (${
        m.ModuleTypeIDName ?? "类型" + moduleList.get(m.ModuleTypeID)
      })`,
      value: m.ModuleUID,
    }));
};
</script>

<style scoped>
.prop-page {
  padding: 8px;
}
</style>
