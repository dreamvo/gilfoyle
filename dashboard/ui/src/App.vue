<template>
  <v-app>
    <v-toolbar flat color="white" light>
      <v-col cols="4" sm="2">
        <RouterLink to="/">
          <v-avatar tile height="32" width="auto">
            <img :src="require('@/assets/logo.svg')" alt="logo" />
          </v-avatar>
        </RouterLink>
      </v-col>

      <v-spacer></v-spacer>

      <v-btn outlined color="#66f">
        <v-icon light>mdi-plus</v-icon>
        Create
      </v-btn>
    </v-toolbar>

    <v-toolbar flat color="#34495e" dark>
      <v-app-bar-nav-icon
        @click="drawerMenu = !drawerMenu"
      ></v-app-bar-nav-icon>

      <v-toolbar-title>Overview</v-toolbar-title>

      <v-spacer></v-spacer>

      <v-system-bar height="54px" dark color="#34495e">
        <v-icon size="8" :color="$store.state.healthy ? 'green' : 'red'"
          >mdi-circle</v-icon
        >
        <span v-if="$store.state.healthy">Instance status: Running</span>
        <span v-else>Instance status: Unavailable</span>
      </v-system-bar>
    </v-toolbar>

    <v-main>
      <v-container>
        <v-row no-gutters>
          <v-col md="2" v-if="drawerMenu">
            <v-list flat>
              <v-subheader>Menu</v-subheader>
              <v-list-item-group color="primary">
                <v-list-item
                  v-for="(item, index) in navigation"
                  :key="index"
                  :to="item.link"
                >
                  <v-list-item-icon>
                    <v-icon v-text="item.icon"></v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title v-text="item.title"></v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-col>

          <v-col :md="drawerMenu ? 10 : 12" style="min-height:90vh;">
            <RouterView />
          </v-col>
        </v-row>
      </v-container>
    </v-main>

    <v-footer dark padless>
      <v-card class="flex" flat tile color="#2e3341">
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
    drawerMenu: true,
    navigation: [
      {
        title: "Overview",
        link: "/",
        icon: "mdi-home"
      },
      {
        title: "Medias",
        link: "/medias",
        icon: "mdi-home"
      },
      {
        title: "Metrics",
        link: "/metrics",
        icon: "mdi-home"
      },
      {
        title: "Settings",
        link: "/settings",
        icon: "mdi-home"
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
  }
});
</script>
