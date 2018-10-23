<template lang="html">
  <div class="wrapper-page">
      <div class="text-center">
          <router-link to = "/" class="logo-lg"><i class="mdi mdi-radar"></i> <span>RuneManager</span></router-link>
      </div>

      <form @submit.prevent = "handleSubmit" class="form-horizontal m-t-20">

          <div class="form-group row">
              <div class="col-12">
                  <div class="input-group">
                      <div class="input-group-prepend">
                          <span class="input-group-text"><i class="mdi mdi-account"></i></span>
                      </div>
                      <input type="text" v-model = "username" name = "username" class="form-control" :class = "{ 'is-invalid': submitted && !username }" placeholder="Username">
                      <div v-show = "submitted && !username" class = "invalid-feedback">Username is required</div>
                  </div>
              </div>
          </div>

          <div class="form-group row">
              <div class="col-12">
                  <div class="input-group">
                      <div class="input-group-prepend">
                          <span class="input-group-text"><i class="mdi mdi-key"></i></span>
                      </div>
                      <input type="password" v-model = "password" name = "password" class="form-control" :class = "{ 'is-invalid': submitted && !password }"  placeholder="Password">
                      <div v-show = "submitted && !username" class = "invalid-feedback">Password is required</div>
                  </div>
              </div>
          </div>

          <div class="form-group row">
              <div class="col-12">
                  <div class="checkbox checkbox-primary">
                      <input id="checkbox-signup" type="checkbox">
                      <label for="checkbox-signup">
                          Remember me
                      </label>
                  </div>

              </div>
          </div>

          <div class="form-group text-right m-t-20">
              <div class="col-xs-12">
                  <button class="btn btn-primary btn-custom w-md waves-effect waves-light" :disabled = "status.loggingIn">Log In</button>
              </div>
          </div>

          <div class="form-group row m-t-30">
              <div class="col-sm-7">
                  <router-link to = "/forgot" class="text-muted"><i class="fa fa-lock m-r-5"></i> Forgot your password?</router-link>
              </div>
              <div class="col-sm-5 text-right">
                  <router-link to="/register" class="text-muted">Create an account</router-link>
              </div>
          </div>
      </form>
  </div>
</template>

<script>
    import { mapState, mapActions } from 'vuex'

    export default {
        data() {
            return {
                username: '',
                password: '',
                submitted: false
            }
        },
        computed: {
            ...mapState('account', ['status'])
        },
        methods: {
            ...mapActions('account', ['login', 'logout']),
            handleSubmit (e) {
                this.submitted = true;
                const { username, password } = this;
                if (this.input.username != "" && this.input.password != "") {
                    this.login({ username, password });
                }
            }
        }
    }

</script>
