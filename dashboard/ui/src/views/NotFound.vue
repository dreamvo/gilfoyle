<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="12">
        <div class="">
          <h1>Page Not Found</h1>
          <p>We're sorry but we can't find the page you're looking for.</p>
          <v-btn outlined color="#66f" to="/">Back To Home</v-btn>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "NotFound",
  data: () => ({})
});
</script>
