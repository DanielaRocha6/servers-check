<template>
  <div>
    <div >
      <loading :active.sync="isLoading"
               :is-full-page="fullPage">
      </loading>
    </div>
    <div v-if="!results">
      <h1>Loading results for http://{{url}}</h1>
    </div>
    <div id="domainResult" v-if="results && !resultsError">
      <b-container >
        <b-row>
          <b-col>
            <h3>{{title}}</h3>
            <img :class="getClass()" :src="logo" @error="replaceByDefault">
            <p>http://{{url}}</p>
          </b-col>`
          <b-col>

            <b-list-group>
              <b-list-group-item>{{isDown ? 'The website is DOWN!' : 'The website is UP!'}}</b-list-group-item>
              <b-list-group-item v-if="sslGrade!==''">The actual lowest ssl grade is '{{sslGrade}}'</b-list-group-item>
              <b-list-group-item v-if="sslGrade===''">There was no ssl grade found</b-list-group-item>
              <b-list-group-item v-if="previousSslGrade!==''">The previous lowest ssl grade is '{{previousSslGrade}}'</b-list-group-item>
              <b-list-group-item v-if="previousSslGrade===''">There's not previous ssl grade</b-list-group-item>
              <b-list-group-item>{{serversChanged ? 'The servers changed in the past hour or more' :'The servers had not changed in the past hour or more' }}</b-list-group-item>
            </b-list-group>

          </b-col>
        </b-row>
      </b-container>

      <br>
      <h4>Servers</h4>
      <b-table striped hover :items="actualServers"></b-table>
    </div>
    <div v-if="results && resultsError">
      <h2><b>Error found</b></h2>
      <p>The domain {{url}} may not exist or there was an error retrieving the data.</p>
      <h5><a :href="`/domains/${url}`"><b>Please try again!</b></a></h5>
    </div>
  </div>
</template>

<script>
import Loading from 'vue-loading-overlay'
import 'vue-loading-overlay/dist/vue-loading.css'
import axios from 'axios'

export default {
  name: 'domains',
  props: {
    url: {
      type: String,
      default: 'Vue!'
    }
  },
  data () {
    return {
      resultsError: false,
      isLoading: true,
      fullPage: true,
      results: false,
      serversChanged: false,
      sslGrade: '',
      previousSslGrade: '',
      title: '',
      logo: '',
      isDown: '',
      actualServers: [],
      logoErr: false
    }
  },
  methods: {
    getDomain () {
      axios({ method: 'GET', url: 'http://localhost:8090/checkDomain/' + this.url, headers: { 'content-type': 'text/plain' } })
        .then(result => {
          var myResult = result.data
          this.actualServers = myResult.Servers
          if (this.actualServers.length !== 0) {
            this.serversChanged = myResult.Servers_Changed
            this.sslGrade = myResult.Ssl_Grade
            this.previousSslGrade = myResult.Previous_ssl_grade
            this.title = myResult.Title
            this.logo = myResult.Logo
            this.isDown = myResult.Is_down
          } else {
            this.resultsError = true
          }
          this.results = true
          this.isLoading = false
        }).catch(error => {
          this.resultsError = true
          this.results = true
          this.isLoading = false
          console.error(error)
        /* eslint-enable */
        })
    },
    replaceByDefault(e) {
      e.target.src = `/not-found.jpg`
      this.logoErr = true
    },
    getClass() {
      if (this.logoErr) {
        return 'errClass'
      }
      return ''
    }
  },
  components: {
    Loading
  },
  mounted () {
    this.getDomain()
  }
}
</script>

<style lang="scss">
  img {
    width: 5rem;
  }
  .errClass {
    width: 13rem;
  }
</style>
