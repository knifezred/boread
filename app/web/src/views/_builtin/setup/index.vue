<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import {
  NButton,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSpace,
} from "naive-ui"
import { fetchSaveDatabaseConfig } from "@/service/api"
import { $t } from "@/locales"

defineOptions({ name: "SetupPage" });

const router = useRouter();

const saving = ref(false);
const saved = ref(false);

const formValue = ref({
  host: "127.0.0.1",
  port: 3306,
  username: "root",
  password: "",
  dbname: "boread",
});

async function handleSave() {
  saving.value = true;
  try {
    const { error } = await fetchSaveDatabaseConfig(formValue.value);
    if (error) {
      window.$message?.error(
        $t("page.setup.saveFailed") + ": " + (error.message || error),
      );
      return;
    }
    saved.value = true;
    window.$message?.success($t("page.setup.saveSuccess"));
  } finally {
    saving.value = false;
  }
}
</script>
