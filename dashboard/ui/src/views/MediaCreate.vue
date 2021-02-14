<template>
  <v-row>
    <v-col cols="12" md="12">
      <MediaForm />
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";
import MediaForm from "../components/MediaForm.vue";

interface Data {}

export default Vue.extend({
  components: { MediaForm },
  data: (): Data => ({}),
  methods: {},
  async created() {}
});
</script>
