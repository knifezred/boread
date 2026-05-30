<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { NButton } from 'naive-ui'

defineOptions({ name: 'BookReader' })

const route = useRoute()
const bookId = (route.params.id || route.query.id) as string
const chapterId = route.params.chapterId as string

const darkMode = ref(false)
const fontSize = ref(18)
const lineHeight = ref(1.8)

const bookInfo = ref({
  id: bookId,
  title: '禁仙！',
  author: '造化斋主',
  cover: '',
  category: '仙侠',
  updateTime: '2026年04月22日',
  wordCount: '17.3万',
  readerCount: '1.02万',
  copyright: '本作品由起点中文网进行电子制作与发行 版权所有 侵权必究'
})

const chapterInfo = ref({
  id: chapterId,
  title: '第一章 天灵根是大白菜',
  updateTime: '2026年04月22日 08:40',
  wordCount: '2764字',
  content: `一间密闭的房间里，一具无头尸体跌坐在床上，从脖颈处喷溅的血液将整张床单染成赤红一片，就连上方的天花板都有不少血迹。

几名刑捕正在勘察现场，收集线索。

"死者是在搬运气血的时候被人用利器一击斩首，死亡时间为昨日晚上八点到九点之间，伤口处并无残留锐劲，无法确认所用武学。"

一名穿着青色捕服，面容冷峻的中年男子正拿着小本子记录手下的汇报，眉头紧锁，显然这个案子有些棘手。

"头，死者身份查清楚了，是青牛镇的张屠户，平日里除了杀猪卖肉也没其他营生，更没什么仇家，好端端的怎么会被人杀了，还砍了脑袋？"

"头，现场没有留下任何脚印和痕迹，凶手明显是个修炼者，而且修为不低，至少是炼气七层以上的修为，不然不可能做到这般悄无声息。"

听着下属的汇报，李捕头的眉头皱的更紧了。

青牛镇只是个边陲小镇，平日里连个炼气五层的修炼者都少见，怎么会突然冒出来一个炼气七层以上的高手，还杀了一个普通的屠户？

事出反常必有妖！

"去，查一下张屠户最近有没有和什么陌生人接触，还有他最近有没有得到过什么东西，事无巨细，都要查清楚！"

"是！"

下属领命离去，李捕头看着地上的尸体，眼神凝重。

就在这时，一道身影从门外走了进来，看到地上的尸体，眉头微微一挑。

"李捕头，什么情况？"

听到声音，李捕头回头一看，脸上的凝重顿时消散，换上了笑容："是林先生啊，您怎么来了？"

来人是个看起来二十岁左右的青年，穿着一身白色的长衫，面容俊朗，气质温润，正是青牛镇唯一的大夫林轩。

"我听说这边出了命案，过来看看有没有什么能帮上忙的。"林轩笑了笑，目光落在地上的尸体上，瞳孔微微一缩，不过很快就恢复了正常。

"劳烦林先生费心了，不过这是我们刑捕房的事，就不劳您费心了。"李捕头笑着说道，虽然语气客气，但是明显不想让林轩插手。

林轩也不在意，点了点头："那行，我就是过来看看，没什么事我就先走了，镇东头的王大爷还等着我去给他看病呢。"

说完，林轩就转身离开了房间，没有人看到，他转身的瞬间，眼底闪过了一丝不易察觉的寒光。

回到自己的医馆，林轩关上房门，脸上的温润消失不见，取而代之的是一片冰冷。

"没想到还是被发现了，看来得尽快离开这里了。"林轩低声呢喃道，手掌一翻，一枚通体漆黑，上面刻满了诡异纹路的戒指出现在手中。

这枚戒指是他昨天从张屠户那里得到的，准确的说，是他杀了张屠户之后抢来的。

他本来也不想杀人，但是张屠户得到这枚戒指的时候，正好被他撞见了，而且张屠户也认出了他的身份，为了不暴露自己的秘密，他只能杀人灭口。

他原本以为自己做的天衣无缝，没想到还是留下了蛛丝马迹，被刑捕房察觉到了不对劲。

"看来这个青牛镇是不能待了，得尽快离开，找个地方突破到炼气八层，不然以我现在的修为，遇到真正的高手根本没有反抗之力。"

林轩心里打定了主意，开始收拾东西，准备连夜离开青牛镇。

就在他收拾东西的时候，突然感觉到一股强大的气息锁定了他，脸色顿时一变。

"不好，被发现了！"

他没有丝毫犹豫，直接破窗而出，朝着镇外的方向掠去。

"林轩，杀了人还想跑？给我站住！"

身后传来了李捕头的怒喝声，林轩回头一看，只见李捕头带着十几个刑捕正朝他追来，李捕头身上的气息赫然已经达到了炼气八层的境界。

"没想到你隐藏的这么深，居然已经是炼气八层的修为了。"林轩眼神凝重，没想到这个看起来普普通通的捕头，居然有这么高的修为。

"少废话，你杀了张屠户，证据确凿，乖乖束手就擒，跟我回刑捕房受审，或许我还能留你一条全尸！"李捕头冷声说道，眼神冰冷。

林轩笑了笑："想抓我？就凭你？恐怕还不够资格！"

话音落下，林轩身上的气息爆发，赫然也是炼气八层的境界，而且比李捕头的气息还要强大几分。

"什么？你居然也是炼气八层！"李捕头脸色大变，显然没想到林轩的修为居然这么高。

"现在知道怕了？晚了！"林轩眼神一冷，直接朝着李捕头冲了过去，一掌拍出，掌风凛冽，带着刺耳的破空声。

李捕头脸色一变，连忙挥拳迎击。

砰！

一声闷响，李捕头整个人直接被林轩一掌拍飞了出去，口中喷出一口鲜血，脸上满是难以置信的神色："怎么可能？你明明和我一样是炼气八层，为什么实力比我强这么多？"

林轩冷笑道："同样是炼气八层，也分三六九等，就你这种靠资源堆上来的废物，也配和我相提并论？"

他从小修炼的就是顶尖的功法，同阶之中几乎无敌，越阶挑战也不是什么难事，更何况李捕头只是个普通的炼气八层，根本不是他的对手。

解决了李捕头，林轩没有丝毫停留，直接朝着远处的山脉掠去，很快就消失在了密林之中。

等到其他刑捕反应过来的时候，林轩已经不见了踪影，只留下躺在地上重伤昏迷的李捕头。

从此，青牛镇再也没有人见过林轩，仿佛这个人从来没有出现过一样。`
})

const bgColor = computed(() => darkMode.value ? '#1a1a1a' : '#f2eede')
const textColor = computed(() => darkMode.value ? '#e0e0e0' : '#333')

const toolbarOptions = [
  { key: 'catalog', label: '目录', icon: 'solar:list-linear' },
  { key: 'review', label: '书评稿', icon: 'solar:document-text-linear' },
  { key: 'bookshelf', label: '加书架', icon: 'solar:bookmark-linear' },
  { key: 'vote', label: '投票', icon: 'solar:like-linear' },
  { key: 'reward', label: '赞赏', icon: 'solar:gift-linear' },
  { key: 'setting', label: '设置', icon: 'solar:settings-linear' },
  { key: 'client', label: '客户端', icon: 'solar:smartphone-linear' },
]

function handleToolbarAction(key: string) {
  console.log('点击工具栏:', key)
}

function prevChapter() {
  console.log('上一章')
}

function nextChapter() {
  console.log('下一章')
}
</script>

<template>
  <div
    class="min-h-screen font-sans relative transition-colors duration-300"
    :class="{ dark: darkMode }"
    :style="{ backgroundColor: bgColor, color: textColor }">
    <header
      class="fixed top-0 left-0 right-0 h-60px bg-white/90 backdrop-blur-md z-100 border-b border-black/5"
      :class="darkMode ? 'bg-[#1a1a1a]/90 border-white/5' : ''">
      <div class="max-w-1200px mx-auto h-full px-6 flex items-center justify-between">
        <div class="flex items-center gap-6">
          <span class="text-lg font-semibold text-[#c4272e]">起点中文网</span>
          <div class="flex items-center gap-1.5 text-xs" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">
            <span>首页</span>
            <SvgIcon icon="solar:chevron-right-linear" size="12" />
            <span>仙侠</span>
            <SvgIcon icon="solar:chevron-right-linear" size="12" />
            <span>修真文明</span>
            <SvgIcon icon="solar:chevron-right-linear" size="12" />
            <span>{{ bookInfo.title }}</span>
          </div>
        </div>
        <div class="flex items-center gap-1.5 text-sm cursor-pointer" :class="darkMode ? 'text-gray-300' : 'text-gray-700'">
          <span class="user-name">大木小棠</span>
          <SvgIcon icon="solar:alt-arrow-down-linear" size="14" />
        </div>
      </div>
    </header>

    <main class="pt-80px pb-100px flex justify-center">
      <div class="w-full max-w-720px px-6">
        <div class="text-center py-10 border-b border-black/8 mb-10" :class="darkMode ? 'border-white/8' : ''">
          <div class="mb-5">
            <div
              class="w-100px h-140px mx-auto rd-1 flex items-center justify-center shadow-lg"
              :style="{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }">
              <span class="text-3xl font-bold text-white">{{ bookInfo.title.slice(0, 2) }}</span>
            </div>
          </div>
          <h1 class="text-3xl font-bold mb-2 text-inherit">{{ bookInfo.title }}</h1>
          <p class="text-sm mb-6" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">{{ bookInfo.author }} 著</p>
          <div class="flex justify-center gap-12 mb-5">
            <div class="text-center">
              <span class="block text-base font-medium mb-1">{{ bookInfo.category }}</span>
              <span class="text-xs" :class="darkMode ? 'text-gray-500' : 'text-gray-400'">文字</span>
            </div>
            <div class="text-center">
              <span class="block text-base font-medium mb-1">{{ bookInfo.updateTime }}</span>
              <span class="text-xs" :class="darkMode ? 'text-gray-500' : 'text-gray-400'">上架</span>
            </div>
            <div class="text-center">
              <span class="block text-base font-medium mb-1">{{ bookInfo.wordCount }}</span>
              <span class="text-xs" :class="darkMode ? 'text-gray-500' : 'text-gray-400'">连载中 (字)</span>
            </div>
          </div>
          <p class="text-sm mb-4" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">与 {{ bookInfo.readerCount }} 位书友共同开启《{{ bookInfo.title }}》的仙侠之旅</p>
          <div class="flex justify-center gap-6 text-xs mb-8" :class="darkMode ? 'text-gray-500' : 'text-gray-400'">
            <span class="flex items-center gap-1">
              <SvgIcon icon="solar:book-open-linear" size="14" />
              总推荐: 987
            </span>
            <span class="flex items-center gap-1">
              <SvgIcon icon="solar:heart-linear" size="14" />
              总收藏: 3647
            </span>
            <span class="flex items-center gap-1">
              <SvgIcon icon="solar:chart-linear" size="14" />
              周点击: 1.8万
            </span>
          </div>
          <div class="flex items-center justify-center gap-1.5 text-xs" :class="darkMode ? 'text-gray-500' : 'text-gray-400'">
            <SvgIcon icon="solar:shield-check-linear" size="14" />
            <span>{{ bookInfo.copyright }}</span>
          </div>
        </div>

        <div class="w-full">
          <div class="mb-10">
            <h2 class="text-2xl font-semibold text-center mb-4 text-inherit">{{ chapterInfo.title }}</h2>
            <div class="flex justify-center gap-6 flex-wrap text-xs" :class="darkMode ? 'text-gray-500' : 'text-gray-400'">
              <span class="flex items-center gap-1">
                <SvgIcon icon="solar:book-linear" size="14" />
                {{ bookInfo.title }}
              </span>
              <span class="flex items-center gap-1">
                <SvgIcon icon="solar:user-linear" size="14" />
                {{ bookInfo.author }}
              </span>
              <span class="flex items-center gap-1">
                <SvgIcon icon="solar:document-text-linear" size="14" />
                {{ chapterInfo.wordCount }}字
              </span>
              <span class="flex items-center gap-1">
                <SvgIcon icon="solar:calendar-linear" size="14" />
                {{ chapterInfo.updateTime }}
              </span>
            </div>
          </div>

          <div class="mb-15 text-inherit" :style="{ fontSize: `${fontSize}px`, lineHeight: lineHeight }">
            <p v-for="(paragraph, index) in chapterInfo.content.split('\n\n')" :key="index" v-if="paragraph.trim()" class="mb-5 text-indent-2em tracking-0.5px">
              {{ paragraph.trim() }}
            </p>
          </div>

          <div class="flex justify-between pt-8 border-t border-black/8" :class="darkMode ? 'border-white/8' : ''">
            <NButton ghost type="primary" size="large" @click="prevChapter">
              <template #icon>
                <SvgIcon icon="solar:arrow-left-linear" size="18" />
              </template>
              上一章
            </NButton>
            <NButton ghost type="primary" size="large" @click="nextChapter">
              下一章
              <template #icon>
                <SvgIcon icon="solar:arrow-right-linear" size="18" />
              </template>
            </NButton>
          </div>
        </div>
      </div>
    </main>

    <aside class="fixed left-6 top-1/2 -translate-y-1/2 flex flex-col gap-5 z-99 max-lg:hidden">
      <div class="w-44px h-44px rd-2 bg-white/80 backdrop-blur-md flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200 shadow-sm text-gray-700 hover:bg-white hover:-translate-y-1" :class="darkMode ? 'bg-[#282828]/80 text-gray-300 hover:bg-[#323232]' : ''" title="反馈">
        <SvgIcon icon="solar:chat-round-dots-linear" size="20" />
      </div>
      <div class="w-44px h-44px rd-2 bg-white/80 backdrop-blur-md flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200 shadow-sm text-gray-700 hover:bg-white hover:-translate-y-1" :class="darkMode ? 'bg-[#282828]/80 text-gray-300 hover:bg-[#323232]' : ''" title="退出">
        <SvgIcon icon="solar:logout-2-linear" size="20" />
      </div>
      <div class="w-44px h-44px rd-2 bg-white/80 backdrop-blur-md flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200 shadow-sm text-gray-700 hover:bg-white hover:-translate-y-1" :class="darkMode ? 'bg-[#282828]/80 text-gray-300 hover:bg-[#323232]' : ''" title="亮度">
        <SvgIcon icon="solar:sun-linear" size="20" />
      </div>
      <div class="w-44px h-44px rd-2 bg-white/80 backdrop-blur-md flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200 shadow-sm text-gray-700 hover:bg-white hover:-translate-y-1" :class="darkMode ? 'bg-[#282828]/80 text-gray-300 hover:bg-[#323232]' : ''" title="字体">
        <SvgIcon icon="solar:text-underline-linear" size="20" />
      </div>
      <div class="w-44px h-44px rd-2 bg-white/80 backdrop-blur-md flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200 shadow-sm text-gray-700 hover:bg-white hover:-translate-y-1" :class="darkMode ? 'bg-[#282828]/80 text-gray-300 hover:bg-[#323232]' : ''" @click="darkMode = !darkMode" title="日间/夜间模式">
        <SvgIcon :icon="darkMode ? 'solar:sun-bold' : 'solar:moon-linear'" size="20" />
      </div>
    </aside>

    <aside class="fixed right-6 top-1/2 -translate-y-1/2 flex flex-col gap-4 z-99 max-lg:hidden">
      <div v-for="item in toolbarOptions" :key="item.key" class="w-44px h-44px rd-2 bg-white/80 backdrop-blur-md flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200 shadow-sm text-gray-700 hover:bg-white hover:-translate-y-1" :class="darkMode ? 'bg-[#282828]/80 text-gray-300 hover:bg-[#323232]' : ''" @click="handleToolbarAction(item.key)" :title="item.label">
        <SvgIcon :icon="item.icon" size="22" />
        <span class="text-[11px] leading-none">{{ item.label }}</span>
      </div>
    </aside>

    <div class="fixed right-30 top-140px cursor-pointer opacity-40 transition-opacity duration-200 hover:opacity-100 text-[#c4272e] max-lg:right-6">
      <SvgIcon icon="solar:bookmark-linear" size="24" />
    </div>
  </div>
</template>

<style scoped>
.text-indent-2em {
  text-indent: 2em;
}
</style>
