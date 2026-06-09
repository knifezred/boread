<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { jsonClone } from "@sa/utils"
import { enableStatusOptions, userGenderOptions } from "@/constants/business"
import {
  fetchGetAllRoles,
  fetchCreateUser,
  fetchUpdateUser,
} from "@/service/api"
import { useFormRules, useNaiveForm } from "@/hooks/common/form"
import { $t } from "@/locales"

defineOptions({
  name: "UserOperateDrawer",
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: Api.SystemManage.User | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: "submitted"): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>("visible", {
  default: false,
});

const { validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const title = computed(() => {
  const titles: Record<NaiveUI.TableOperateType, string> = {
    add: $t("page.admin.system.user.addUser"),
    edit: $t("page.admin.system.user.editUser"),
  };
  return titles[props.operateType];
});

type Model = Pick<
  Api.SystemManage.User,
  | "id"
  | "userName"
  | "userGender"
  | "nickName"
  | "userPhone"
  | "userEmail"
  | "userRoles"
  | "status"
  | "password"
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    userName: "",
    password: "",
    userGender: null,
    nickName: "",
    userPhone: "",
    userEmail: "",
    userRoles: [],
    status: "1",
  };
}

type RuleKey = Extract<keyof Model, "userName" | "status" | "password">;

const rules: Record<RuleKey, App.Global.FormRule> = {
  userName: defaultRequiredRule,
  status: defaultRequiredRule,
  password: {
    required: props.operateType === "add",
    min: 6,
    max: 64,
    message: $t("page.admin.system.user.form.password"),
    trigger: "input",
  },
};

/** the enabled role options (value = numeric id, 与后端 roleIds 对齐) */
const roleOptions = ref<CommonType.Option<number>[]>([]);

/** 当前选中的角色 ID 列表 */
const selectedRoleIds = ref<number[]>([]);

async function getRoleOptions() {
  const { error, data } = await fetchGetAllRoles();

  if (!error) {
    roleOptions.value = data.map((item) => ({
      label: item.roleName,
      value: item.id,
    }));
  }
}

function handleInitModel() {
  model.value = createDefaultModel();
  selectedRoleIds.value = [];

  if (props.operateType === "edit" && props.rowData) {
    Object.assign(model.value, jsonClone(props.rowData));
    // 后端直接返回 roleIds，无需转换
    selectedRoleIds.value = props.rowData.roleIds ?? [];
  }
}

function closeDrawer() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();
  // request
  const payload = {
    ...model.value,
    roleIds: selectedRoleIds.value,
  };
  let requestResult;
  if (props.operateType === "edit") {
    requestResult = await fetchUpdateUser(model.value.id, payload);
  } else {
    requestResult = await fetchCreateUser(payload);
  }
  if (requestResult.error) {
    window.$message?.error(requestResult.error.message);
  }
  closeDrawer();
  emit("submitted");
}

watch(visible, () => {
  if (visible.value) {
    restoreValidation();
    handleInitModel();
    getRoleOptions();
  }
});
</script>

<template>
  <NDrawer v-model:show="visible" display-directive="show" :width="360">
    <NDrawerContent :title="title" :native-scrollbar="false" closable>
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem
          :label="$t('page.admin.system.user.userName')"
          path="userName"
        >
          <NInput
            v-model:value="model.userName"
            :placeholder="$t('page.admin.system.user.form.userName')"
          />
        </NFormItem>
        <NFormItem
          v-if="operateType === 'add'"
          :label="$t('page.admin.system.user.password')"
          path="password"
        >
          <NInput
            v-model:value="model.password"
            type="password"
            show-password-on="click"
            :placeholder="$t('page.admin.system.user.form.password')"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.system.user.userGender')"
          path="userGender"
        >
          <NRadioGroup v-model:value="model.userGender">
            <NRadio
              v-for="item in userGenderOptions"
              :key="item.value"
              :value="item.value"
              :label="$t(item.label)"
            />
          </NRadioGroup>
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.system.user.nickName')"
          path="nickName"
        >
          <NInput
            v-model:value="model.nickName"
            :placeholder="$t('page.admin.system.user.form.nickName')"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.system.user.userPhone')"
          path="userPhone"
        >
          <NInput
            v-model:value="model.userPhone"
            :placeholder="$t('page.admin.system.user.form.userPhone')"
          />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.user.userEmail')" path="email">
          <NInput
            v-model:value="model.userEmail"
            :placeholder="$t('page.admin.system.user.form.userEmail')"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.system.user.userStatus')"
          path="status"
        >
          <NRadioGroup v-model:value="model.status">
            <NRadio
              v-for="item in enableStatusOptions"
              :key="item.value"
              :value="item.value"
              :label="$t(item.label)"
            />
          </NRadioGroup>
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.user.userRole')" path="roles">
          <NSelect
            v-model:value="selectedRoleIds"
            multiple
            :options="roleOptions"
            :placeholder="$t('page.admin.system.user.form.userRole')"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace :size="16">
          <NButton @click="closeDrawer">{{ $t("common.cancel") }}</NButton>
          <NButton type="primary" @click="handleSubmit">
            {{
              $t("common.confirm")
            }}
          </NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
