<template>
  <div>
    <div class="prop-page" v-if="store.selectedAction">
      <h3 style="margin-bottom: 16px; color: black">
        Action 编辑
        <a-form layout="vertical">
          <a-form-item label="Action名">
            <a-input v-model:value="store.selectedAction.Name" />
          </a-form-item>
        </a-form>
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
          <a-collapse v-model:activeKey="ioSubModuleActiveKey">
            <a-collapse-panel
              v-for="(mod, index) in store.selectedAction.TypeFeatureField
                .Modules"
              :key="mod.ModuleUID"
              :header="moduleList.get(mod.ModuleTypeID)"
            >
              <a-divider
                >子模块 {{ moduleList.get(mod.ModuleTypeID) }}</a-divider
              >

              <a-form-item label="ModuleUID">
                <a-input-number
                  v-model:value="mod.ModuleUID"
                  disabled
                  style="width: 100%"
                />
              </a-form-item>

              <!-- Fill -->
              <template v-if="mod.ModuleTypeID === 10">
                <a-form-item label="使用变量">
                  <a-input v-model:value="mod.UseVar" />
                  <a-button type="link" size="small" @click="showVarHelp">
                    查看说明
                  </a-button>
                </a-form-item>
              </template>

              <!-- Fixed -->
              <template v-else-if="mod.ModuleTypeID === 11">
                <a-form-item label="固定内容 (字节数组)" tooltip="填写 -1 跳过内容检查 (发送不可用)">
                  <a-textarea
                    :value="fixedToString(mod.FixedContent)"
                    @blur="(e : any)=> mod.FixedContent = fixedFromString(e.target.value)"
                  />
                </a-form-item>
              </template>

              <!-- Calc -->
              <template v-else-if="mod.ModuleTypeID === 12">
                <a-form-item label="模式">
                  <a-select
                    v-model:value="mod.Mode"
                    style="width: 100%"
                    :options="CalcMode"
                    placeholder="请选择计算模式"
                  />
                </a-form-item>
                <a-form-item label="计算函数">
                  <a-select
                    v-model:value="mod.CalcFunc"
                    style="width: 100%"
                    :options="calcFn"
                    placeholder="请选择计算函数"
                  />
                </a-form-item>
                <a-form-item
                  label="计算时机"
                  tooltip="如校验和计算需要完整的数据，则需要选择组装后；长度计算则不需要(带有Fill模块时除外)，数据长度开始时已经计算完毕"
                >
                  <a-select
                    v-model:value="mod.CalcTiming"
                    style="width: 100%"
                    :options="CalcTiming"
                    placeholder="请选择计算时机"
                  />
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
                <a-form-item
                  label="自定义内容 (字节数组)"
                  tooltip="填写 -1 跳过内容检查"
                >
                  <a-textarea
                    :value="fixedToString(mod.CustomContent)"
                    @blur="
                  (e : any)=> mod.CustomContent = fixedFromString(e.target.value)
                "
                  />
                </a-form-item>
                <!-- 多选输入模块UID -->
                <a-form-item
                  label="参考输入长度模块 UID"
                  :rules="[{ required: false }]"
                  tooltip="指定模块接收到的数据作为custom模块的读取长度"
                  v-if="store.selectedAction.ActionTypeID == 2"
                >
                  <a-select
                    v-model:value="mod.ReceiveVarLengthModuleUID"
                    style="width: 100%"
                    :options="calcInputOptions(mod)"
                    placeholder="请选择输入模块"
                    allowClear
                  />
                </a-form-item>
              </template>
              <a-button danger @click="removeModule(index)">删除模块</a-button>
              <Divider></Divider>
            </a-collapse-panel>
          </a-collapse>
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
            <!-- 变量名 -->
            <a-form-item label="变量名">
              <a-input
                v-model:value="store.selectedAction.TypeFeatureField.VarName"
              />
            </a-form-item>

            <!-- 变量类型选择 -->
            <a-form-item label="变量类型">
              <a-select
                v-model:value="store.selectedAction.TypeFeatureField.VarType"
                :options="varTypeOptions"
                placeholder="请选择变量类型"
                style="width: 100%"
              />
            </a-form-item>

            <!-- 根据变量类型显示不同的输入框 -->
            <a-form-item label="变量值">
              <!-- Number 类型 -->
              <a-input-number
                v-if="
                  store.selectedAction.TypeFeatureField.VarType === 'number'
                "
                v-model:value="
                  store.selectedAction.TypeFeatureField.VarNumberValue
                "
                style="width: 100%"
              />

              <!-- String 类型 -->
              <a-input
                v-else-if="
                  store.selectedAction.TypeFeatureField.VarType === 'string'
                "
                v-model:value="
                  store.selectedAction.TypeFeatureField.VarStringValue
                "
              />

              <!-- ByteArray 类型 -->
              <a-textarea
                v-else-if="
                  store.selectedAction?.TypeFeatureField?.VarType === 'array'
                "
                :value="
                  fixedToString(
                    store.selectedAction?.TypeFeatureField?.VarByteArrayValue
                  )
                "
                @blur="(e: any) => {
            if (store.selectedAction && store.selectedAction.TypeFeatureField) {
              store.selectedAction.TypeFeatureField.VarByteArrayValue = fixedFromString(e.target.value)
            }
          }"
              />

              <!-- JSON 类型（vue-json-editor） -->
              <json-editor-vue
                v-else-if="
                  store.selectedAction.TypeFeatureField.VarType === 'JSON'
                "
                v-model="store.selectedAction.TypeFeatureField.VarStringValue"
                :mode="Mode.text"
                lang="zh"
                height="400px"
              />

              <!-- 如果没选类型 -->
              <span v-else style="color: #999">请选择变量类型</span>
            </a-form-item>
          </template>

          <!-- Assign -->
          <template v-else-if="store.selectedAction.ActionTypeID === 32">
            <a-form-item label="赋值变量">
              <a-input
                v-model:value="store.selectedAction.TypeFeatureField.AssignTargetVar"
              />
            </a-form-item>
            <a-form-item label="赋值表达式">
              <a-input
                v-model:value="store.selectedAction.TypeFeatureField.Expression"
              />
              <a-button type="link" size="small" @click="showFmtHelp">
                查看说明
              </a-button>
            </a-form-item>
          </template>

          <!-- IF -->
          <template v-else-if="store.selectedAction.ActionTypeID === 24">
            <a-form-item label="条件表达式">
              <a-input
                v-model:value="store.selectedAction.TypeFeatureField.Condition"
              />
              <a-button type="link" size="small" @click="showFmtHelp">
                查看说明
              </a-button>
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
              <a-button type="link" size="small" @click="showFmtHelp">
                查看说明
              </a-button>
            </a-form-item>
          </template>

          <!-- Show -->
          <template v-if="store.selectedAction.ActionTypeID === 92">
            <a-form-item label="输出目标下标 (按文字框绑定顺序 index:[0-1])">
              <a-select
                v-model:value="store.selectedAction.TypeFeatureField.OutputIdx"
                :options="outputIdxList"
                placeholder="请选择目标文字框index"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item label="打印格式字符串">
              <a-input
                v-model:value="store.selectedAction.TypeFeatureField.FmtStr"
              />
              <a-button type="link" size="small" @click="showFmtHelp">
                查看说明
              </a-button>
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

    <a-modal
      v-model:open="helpFmtVisible"
      title="格式化字符串\表达式 说明"
      footer=""
      width="600px"
    >
      <p>可以使用占位符获取指定内容：</p>
      <ul>
        <li><code>{0}</code> 上一Action名称</li>
        <li><code>{1}</code> 上一串口接收数据</li>
        <li><code>{2}</code> 上一Action执行结果</li>
        <li><code>{3}</code> 当前时间</li>
        <li><code>{foo}</code> foo变量的值</li>
      </ul>

      <p>表达式支持以下语法：</p>
      <ul>
        <li>&& || < <= > >= !=</li>
      </ul>

      <p>数组下标语法：</p>
      <ul>
        <li><code>{goo:0}</code> goo数组的第一个字节 等价于goo[0]</li>
        <li><code>{goo:0,3}</code> goo数组的[0,3)字节 等价于goo[0:3]</li>
      </ul>
    </a-modal>

    <a-modal
      v-model:open="helpVarVisible"
      title="变量 说明"
      footer=""
      width="600px"
    >
      <p>数组下标语法：</p>
      <ul>
        <li><code>goo:0</code> goo数组的第一个字节 等价于goo[0]</li>
        <li><code>goo:0,3</code> goo数组的[0,3)字节 等价于goo[0:3]</li>
      </ul>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, watch, ref, onMounted } from "vue";
import { useActionStore } from "../../stores/action_store";
import { parseActionTags } from "../../utils/action_utils";
import { GetAllCalcFn } from "../../../wailsjs/go/bsd_testtool/Manager";
import { Divider } from "ant-design-vue";
import JsonEditorVue from "json-editor-vue";
import { Mode } from "vanilla-jsoneditor";

const store = useActionStore();

onMounted(() => {
  getCalcFn();
});

const newModuleType = ref<number | null>(null);

const helpFmtVisible = ref(false);
const helpVarVisible = ref(false);
const showFmtHelp = () => (helpFmtVisible.value = true);
const showVarHelp = () => (helpVarVisible.value = true);

const ioSubModuleActiveKey = ref([]);

const varTypeOptions = [
  { label: "数字(Number)", value: "number" },
  { label: "文本(String)", value: "string" },
  { label: "字节数组(ByteArray)", value: "array" },
  { label: "JSON", value: "JSON" },
];

const outputIdxList = [{label: '0', value:0}, {label: '1', value:1}]

watch(
  () => store.selectedAction?.TypeFeatureField,
  (newVal) => {
    if (store.selectedAction) {
      const sa = store.selectedAction as any;
      sa.Tags = parseActionTags(store.selectedAction);
      store.selectedAction = sa;
      console.log("change ");
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
        ReceiveVarLengthModuleUID: typeId,
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
      return `${hex}`; // 加前缀，便于识别
    })
    .join(",");
};
const fixedFromString = (str: any) => {
  console.log(str);

  if (!str) return [];
  return str
    .split(",")
    .map((s: any) => parseInt(s, 16))
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

const CalcMode = [
  { label: "计算", value: "Calc" },
  { label: "检查", value: "Check" },
];

const CalcTiming = [
  { label: "即时", value: "Now" },
  { label: "组装后", value: "Post" },
];

const calcFn = ref<{ label: string; value: string }[]>();

const getCalcFn = () => {
  GetAllCalcFn().then((res) => {
    const fnList = res.map((v: string) => ({
      label: v,
      value: v,
    }));
    calcFn.value = fnList;
  });
};

const calcInputOptions = (currentMod: any) => {
  const modules = store.selectedAction?.TypeFeatureField?.Modules ?? [];
  return modules
    .filter((m: any) => m.ModuleUID !== currentMod.ModuleUID) // 排除自己
    .map((m: any) => ({
      label: `${
        m.ModuleTypeIDName ?? "类型" + moduleList.get(m.ModuleTypeID)
      } (${m.ModuleUID})`,
      value: m.ModuleUID,
    }));
};
</script>

<style scoped>
.prop-page {
  padding: 8px;
}
</style>
