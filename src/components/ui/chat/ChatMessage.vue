<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import markdownit from 'markdown-it';
import hljs from 'highlight.js';
import { computed, ref, watch } from 'vue';
import 'highlight.js/styles/github.css';

const { t } = useI18n();

const props = defineProps<{
  message: ChatMessage;
}>();

// Define interfaces for content items
interface BaseContentItem {
  type: string;
  content: string;
  isStreaming?: boolean;
}

interface TextContentItem extends BaseContentItem {
  type: 'text';
}

interface ActionContentItem extends BaseContentItem {
  type: 'action';
  action: string;
  action_status: string;
  action_id?: string;
}

type ContentItem = TextContentItem | ActionContentItem;

const md = markdownit({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        // Add data-language attribute to pre tag
        return (
          '<pre data-language="' +
          lang +
          '"><code class="hljs language-' +
          lang +
          '">' +
          hljs.highlight(str, { language: lang }).value +
          '</code></pre>'
        );
      } catch (__) {}
    }
    // For unknown languages, still add the language if provided
    return lang
      ? '<pre data-language="' + lang + '"><code>' + str + '</code></pre>'
      : '';
  },
});

// Format timestamp
const formatTime = (date: Date | undefined): string => {
  return date
    ? date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    : '';
};

// Safe markdown rendering with sanitization
const renderMarkdown = (content: string): string => {
  if (!content) return '';
  return md.render(content);
};

const typing = ref(true);

// Stop the typing animation when streaming ends
watch(
  () => props.message.isStreaming,
  isStreaming => {
    typing.value = isStreaming === true;
  },
  { immediate: true }
);

// Function to get action status icon
const getActionStatusIcon = (status: string) => {
  switch (status) {
    case 'pending':
      return ['fas', 'circle-notch'];
    case 'success':
      return ['fas', 'check-circle'];
    case 'failed':
      return ['fas', 'times-circle'];
    default:
      return ['fas', 'circle-notch'];
  }
};

// Function to get action status color
const getActionStatusColor = (status: string) => {
  switch (status) {
    case 'pending':
      return 'var(--el-color-warning)';
    case 'success':
      return 'var(--el-color-success)';
    case 'failed':
      return 'var(--el-color-danger)';
    default:
      return 'var(--el-color-info)';
  }
};

// Function to format action name for display
const formatActionName = (action: string) => {
  return action
    .replace(/[_:]/g, ' ')
    .replace(/\b\w/g, char => char.toUpperCase());
};

// Compute content items from message
const contentItems = computed((): ContentItem[] => {
  const items: ContentItem[] = [];

  // First check for structured content items
  if (props.message.contents && props.message.contents.length > 0) {
    // Use the pre-structured content items from the message
    props.message.contents.forEach(content => {
      if (content.type === 'action') {
        items.push({
          type: 'action',
          content: content.content || '',
          action: content.action || '',
          action_status: content.action_status || 'pending',
          action_id: content.key,
          isStreaming: false,
        } as ActionContentItem);
      } else {
        items.push({
          type: 'text',
          content: content.content || '',
          isStreaming: false,
        } as TextContentItem);
      }
    });
    return items;
  }

  // Add the main text content if it exists
  if (props.message.content) {
    items.push({
      type: 'text',
      content: props.message.content,
      isStreaming: props.message.isStreaming,
    } as TextContentItem);
  }

  return items;
});

defineOptions({ name: 'ClChatMessage' });
</script>

<template>
  <div :class="['message-container', message.role]">
    <div class="message-content">
      <template v-if="message.content">
        <div v-html="renderMarkdown(message.content)"></div>
      </template>
      <!-- Iterate through content items in order -->
      <div v-else class="content-items">
        <template v-for="(content, index) in message.contents" :key="index">
          <!-- Action content -->
          <div
            v-if="content.type === 'action'"
            class="action-item"
            :class="`action-status-${content.action_status}`"
          >
            <div class="action-header">
              <cl-icon
                :icon="getActionStatusIcon(content.action_status!)"
                :style="{
                  color: getActionStatusColor(content.action_status!),
                }"
                class="action-icon"
                :class="{
                  spinning: content.action_status === 'pending',
                }"
              />
              <span class="action-name">{{
                formatActionName(content.action!)
              }}</span>
              <span
                class="action-status"
                :style="{
                  color: getActionStatusColor(content.action_status!),
                }"
              >
                {{ content.action_status }}
              </span>
            </div>
            <div v-if="content.content" class="action-content">
              {{ content.content }}
            </div>
          </div>

          <!-- Text content -->
          <div v-else-if="content.type === 'text'" class="text-content">
            <template v-if="content.isStreaming">
              <div v-html="renderMarkdown(content.content || '')"></div>
              <span class="typing-indicator" v-if="typing">|</span>
            </template>
            <template v-else>
              <div v-html="renderMarkdown(content.content || '')"></div>
            </template>
          </div>
        </template>
      </div>
    </div>

    <div class="message-time">
      <!-- Show 'Generating...' for streaming messages -->
      <template v-if="message.isStreaming">
        <span class="typing-text">{{
          t('components.ai.chatbot.generating')
        }}</span>
      </template>
      <template v-else>
        {{ formatTime(message.timestamp) }}
      </template>
    </div>
  </div>
</template>

<style scoped>
.message-container {
  font-size: 14px;
  width: calc(100% - 24px);
  margin: 0 12px;
  padding: 12px;
  position: relative;
}

.message-container:first-child.user {
  margin-top: 12px;
}

.message-container.user {
  border-radius: 12px;
  background-color: var(--el-color-primary-dark-2);
}

.message-container.system {
  background-color: transparent;
}

.message-content {
  word-break: break-word;
  line-height: 1.5;
  max-height: 100%;
  overflow-y: auto;

  /* Firefox scrollbar styles */
  scrollbar-width: thin;
  scrollbar-color: var(--el-border-color-darker) var(--el-fill-color-lighter);

  /* Webkit scrollbar styles */

  &::-webkit-scrollbar {
    width: 7px;
    background-color: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background-color: var(--el-border-color-darker);
    border-radius: 3px;
  }

  &::-webkit-scrollbar-track {
    background-color: var(--el-fill-color-lighter);
    border-radius: 3px;
  }
}

.content-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.text-content {
  width: 100%;
}

/* Add styles for markdown elements */
.message-content :deep(pre) {
  background-color: var(--el-fill-color-light);
  padding: 24px 12px 12px 12px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 12px 0;
  position: relative;
}

.message-content :deep(pre code) {
  /* Firefox scrollbar styles */
  scrollbar-width: thin;
  scrollbar-color: var(--el-border-color-darker) var(--el-fill-color-lighter);

  /* Webkit scrollbar styles */

  &::-webkit-scrollbar {
    height: 7px;
    background-color: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background-color: var(--el-border-color-darker);
    border-radius: 3px;
  }

  &::-webkit-scrollbar-track {
    background-color: var(--el-fill-color-lighter);
    border-radius: 3px;
  }
}

/* Add language display */
.message-content :deep(pre[data-language]::before) {
  content: attr(data-language);
  position: absolute;
  top: 0;
  left: 0;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  background: var(--el-fill-color);
  padding: 2px 8px;
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
  text-transform: lowercase;
}

.message-content :deep(pre code) {
  font-family: 'Menlo', 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  display: block;
  padding: 0;
  background: none;
}

.message-content :deep(code:not(pre code)) {
  background-color: var(--el-fill-color-light);
  padding: 2px 4px;
  border-radius: 4px;
  font-size: 0.9em;
}

.message-content :deep(a) {
  color: var(--el-color-primary);
  text-decoration: none;
}

.message-content :deep(blockquote) {
  border-left: 4px solid var(--el-color-info-light-5);
  padding-left: 12px;
  margin-left: 0;
  color: var(--el-text-color-secondary);
}

.message-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin-bottom: 16px;
}

.message-content :deep(th),
.message-content :deep(td) {
  border: 1px solid var(--el-border-color);
  padding: 8px;
  text-align: left;
}

.message-content :deep(th) {
  background-color: var(--el-color-info-light-9);
}

.message-content :deep(p) {
  margin: 0;
  padding: 0;
}

.message-content :deep(p:not(:first-child)) {
  margin-top: 6px;
}

.message-content :deep(h1) {
  margin: 0;
  padding: 12px 0;
}

.message-content :deep(h2) {
  margin: 0;
  padding: 12px 0;
}

.message-content :deep(h3) {
  margin: 0;
  padding: 12px 0;
}

.message-content :deep(h4) {
  margin: 0;
  padding: 12px 0;
}

.message-content :deep(h5) {
  margin: 0;
  padding: 12px 0;
}

.message-content :deep(h6) {
  margin: 0;
  padding: 12px 0;
}

.message-content :deep(ul) {
  margin: 0;
  padding-inline-start: 24px;
}

.message-content :deep(ol) {
  margin: 0;
  padding-inline-start: 24px;
}

.message-content :deep(li) {
  margin: 3px 0;
}

.message-time {
  font-size: 10px;
  opacity: 0.7;
  margin-top: 6px;
}

.message-container.user .message-content {
  color: var(--el-color-white);
}

.message-container.user .message-time {
  color: var(--el-color-white);
}

.message-container.system .message-content {
  color: var(--el-text-color-regular);
}

.message-container.system .message-time {
  color: var(--el-text-color-regular);
}

.typing-indicator {
  display: inline-block;
  animation: blink 1s infinite;
  margin-left: 2px;
}

.typing-text {
  display: inline-block;
  color: var(--el-color-primary);
}

@keyframes blink {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0;
  }
}

/* Action styles */
.action-item {
  padding: 8px 12px;
  border-radius: 6px;
  background-color: var(--el-fill-color-light);
  border-left: 3px solid var(--el-color-info);
}

.action-item.action-status-pending {
  border-left-color: var(--el-color-warning);
}

.action-item.action-status-success {
  border-left-color: var(--el-color-success);
}

.action-item.action-status-failed {
  border-left-color: var(--el-color-danger);
}

.action-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.action-icon {
  font-size: 14px;
}

.spinning {
  animation: spin 1s linear infinite;
}

.action-name {
  flex: 1;
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.action-status {
  font-size: 12px;
  text-transform: capitalize;
}

.action-content {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin-left: 24px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
