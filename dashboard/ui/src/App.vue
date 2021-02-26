<template>
  <v-app>
    <v-app-bar elevation="0" clipped-left fixed app light>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      <v-toolbar-title>
        <router-link to="/">
          <v-avatar tile height="32" width="auto">
            <img :src="require('@/assets/logo.svg')" alt="logo" />
          </v-avatar>
        </router-link>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <v-system-bar height="54px" color="transparent">
        <v-icon size="8" :color="$store.state.healthy ? 'green' : 'red'"
          >mdi-circle
        </v-icon>
        <span v-if="$store.state.healthy">Instance status: Running</span>
        <span v-else>Instance status: Unavailable</span>
      </v-system-bar>

      <v-btn small depressed dark color="primary" class="ml-3" to="/medias/create">
        <v-icon light>mdi-plus</v-icon>
        Create media
      </v-btn>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" clipped fixed app light color="white">
      <v-list>
        <v-list-item
          v-for="(item, i) in navigation"
          :key="i"
          :to="item.to"
          router
          exact
        >
          <v-list-item-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title v-text="item.title" />
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-content>
      <v-container fluid>
        <v-row v-if="loading">
          <v-col cols="12" md="6">
            <v-skeleton-loader
              type="card-heading, list-item-three-line, actions"
            ></v-skeleton-loader>
            <v-skeleton-loader
              type="article, card, actions"
            ></v-skeleton-loader>
          </v-col>
          <v-col cols="12" md="6">
            <v-skeleton-loader
              type="date-picker, list-item-three-line"
            ></v-skeleton-loader>
          </v-col>
        </v-row>
        <v-container v-else>
          <RouterView v-if="$store.state.healthy" />
          <v-row v-else>
            <v-col cols="12" md="12">
              <h1>Unhealthy instance</h1>
              <p>
                Something bad is happening, you should look at the logs of your
                Gilfoyle instance.
              </p>
              <h3>Diagnostic</h3>
              <p>{{ $store.state.healthError }}</p>
            </v-col>
          </v-row>
        </v-container>
      </v-container>
    </v-content>

    <v-footer app dark padless>
      <v-card class="flex" flat tile color="secondary">
        <v-card-title>
          <strong class="subheading">Get connected with us!</strong>

          <v-spacer></v-spacer>

          <v-btn
            v-for="(icon, index) in footerIcons"
            :key="index"
            :href="icon.link"
            target="_blank"
            class="mx-4"
            dark
            icon
          >
            <v-icon size="24">
              {{ icon.icon }}
            </v-icon>
          </v-btn>
        </v-card-title>
      </v-card>
    </v-footer>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import config from "./config";
import store from "./store";

export default Vue.extend({
  name: "App",
  components: {},
  data: () => ({
    drawer: true,
    loading: true,
    navigation: [
      {
        title: "Overview",
        to: "/",
        icon: "mdi-view-dashboard"
      },
      {
        title: "Medias",
        to: "/medias",
        icon: "mdi-animation-play"
      },
      {
        title: "About",
        to: "/about",
        icon: "mdi-information"
      }
    ],
    footerIcons: [
      {
        icon: "mdi-github",
        link: config.links.githubURL
      },
      {
        icon: "mdi-twitter",
        link: config.links.twitterURL
      },
      {
        icon: "mdi-web",
        link: config.links.websiteURL
      }
    ]
  }),
  methods: {},
  async created() {
    await store.dispatch("healthCheck");

    setInterval(async () => {
      await store.dispatch("healthCheck");
    }, config.healthCheckDelaySeconds * 1000);

    this.loading = false;
  }
});
</script>
