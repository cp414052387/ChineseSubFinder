<template>
  <div style="min-height: 300px">
    <q-list v-if="csfSearchResult?.length" separator>
      <q-item v-for="(item, index) in csfSearchResult" :key="item.sub_sha256">
        <q-item-section>
          <div class="row items-center q-gutter-sm">
            <div>{{ index + 1 }}. {{ item.title }}</div>
            <q-badge color="primary">{{ LANGUAGES[item.language] }}</q-badge>
            <q-badge
              v-if="cacheBlob[item.sub_sha256]"
              color="grey"
              title="已缓存到浏览器，再次预览和下载不消耗次数，关闭窗口后失效"
              >已缓存</q-badge
            >
          </div>
        </q-item-section>
        <q-item-section side>
          <div class="row">
            <btn-dialog-preview-video
              :path="path"
              :subtitle-url-list="[selectedSubUrl]"
              :on-btn-click="(callback) => handlePreviewClick(item, callback)"
              :subtitle-type="selectedItem?.ext.replace('.', '')"
            />
            <q-btn color="primary" icon="download" flat dense @click="handleDownloadCsfSub(item)" title="下载" />
          </div>
        </q-item-section>
      </q-item>
    </q-list>
    <div v-else-if="!loading" class="text-grey">
      <template v-if="assrtApiErrorMsg">
        <div class="text-negative">获取字幕列表失败，错误信息：{{ assrtApiErrorMsg }}</div>
        <div><q-btn flat label="重试" color="primary" dense @click="searchAssrt" /></div>
      </template>
      <template v-else>
        <div>未搜索到数据，<q-btn flat label="重试" color="primary" dense @click="searchAssrt" /></div>
        <div>如果报错信息提示没有 Token，请到<b>配置中心-字幕源设置</b>，填写Assrt的Token</div>
      </template>
    </div>
    <q-inner-loading :showing="loading">
      <q-spinner size="50px" color="primary" />
      <div>{{ loadingMsg }}</div>
    </q-inner-loading>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import LibraryApi from 'src/api/LibraryApi';
import { SystemMessage } from 'src/utils/message';
import AssrtSubtitlesApi from 'src/api/AssrtSubtitlesApi';
import BtnDialogPreviewVideo from 'pages/library/BtnDialogPreviewVideo.vue';
import { getSubtitleUploadList } from 'pages/library/use-library';
import eventBus from 'vue3-eventbus';
import { useQuasar } from 'quasar';
import { LANGUAGES } from 'src/constants/LibraryConstants';
import { VIDEO_TYPE_MOVIE, VIDEO_TYPE_TV } from 'src/constants/SettingConstants';

const props = defineProps({
  path: String,
  isMovie: {
    type: Boolean,
    default: false,
  },
  season: {
    type: Number,
  },
  episode: {
    type: Number,
  },
});

const $q = useQuasar();
const assrtApiErrorMsg = ref('');
const loading = ref(false);
const loadingMsg = ref('');
const csfSearchResult = ref(null);
const selectedSubBlob = ref(null);
const selectedItem = ref(null);

// blob缓存
const cacheBlob = ref({});

const selectedSubUrl = computed(() => {
  if (selectedSubBlob.value) {
    return URL.createObjectURL(selectedSubBlob.value);
  }
  return null;
});

const setLock = async () => {
  const [, err] = await LibraryApi.setSkipInfo({
    video_skip_infos: [
      {
        video_type: props.isMovie ? VIDEO_TYPE_MOVIE : VIDEO_TYPE_TV,
        physical_video_file_full_path: props.path,
        is_bluray: false,
        is_skip: true,
      },
    ],
  });
  if (err !== null) {
    SystemMessage.error(err.message);
  } else {
    // 通知列表页锁定成功
    eventBus.emit(`refresh-skip-status-${props.path}`, true);
    SystemMessage.success('操作成功');
  }
};

const searchAssrt = async () => {
  loading.value = true;
  assrtApiErrorMsg.value = '';
  loadingMsg.value = '正在获取字幕列表...';
  const [data, err] = await AssrtSubtitlesApi.search({
    is_movie: props.isMovie,
    video_f_path: props.path,
  });
  if (err !== null) {
    assrtApiErrorMsg.value = err.message;
    SystemMessage.error(err.message);
  } else {
    csfSearchResult.value = data.subtitles;
  }
  loadingMsg.value = '';
  loading.value = false;
};

const fetchSubtitleBlob = async (item) => {
  selectedItem.value = item;
  if (cacheBlob.value[item.sub_sha256]) {
    selectedSubBlob.value = cacheBlob.value[item.sub_sha256];
    return;
  }
  selectedSubBlob.value = null;
  loading.value = true;
  loadingMsg.value = '正在获取下载地址...';
  loadingMsg.value = '正在下载字幕...';
  const [blob, err] = await AssrtSubtitlesApi.download(item.sub_sha256);
  if (err !== null) {
    SystemMessage.error(err.message);
  } else {
    cacheBlob.value = {
      ...cacheBlob.value,
      [item.sub_sha256]: blob,
    };
    selectedSubBlob.value = blob;
  }
  loadingMsg.value = '';
  loading.value = false;
};

const handleDownloadCsfSub = async (item) => {
  await fetchSubtitleBlob(item);

  if (!selectedSubBlob.value) {
    return;
  }

  // 上传
  const formData = new FormData();
  formData.append('video_f_path', props.path);
  formData.append('file', new File([selectedSubBlob.value], item.title, { type: 'text/plain' }));
  await LibraryApi.uploadSubtitle(formData);
  await getSubtitleUploadList();
  eventBus.emit('subtitle-uploaded');

  $q.dialog({
    title: '操作确认',
    message: `已下载到库中，是否锁定该视频，无需再次自动下载字幕？`,
    cancel: true,
    persistent: true,
    focus: 'none',
  }).onOk(async () => {
    setLock();
  });

  SystemMessage.success('已下载到库中');
};

const handlePreviewClick = async (item, callback) => {
  await fetchSubtitleBlob(item);
  if (selectedSubUrl.value) {
    callback(true);
  } else {
    callback(false);
  }
};

onMounted(() => {
  searchAssrt();
});
</script>
