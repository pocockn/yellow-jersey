### TODO:

INFRA
- Add some CI
- Release process 
- Add linter (yawn)
- K8s set up - Minikube
- Add the mono repo into Git ✔️
- Once release process is finalised implement GoReleaser

CONFIG/SECURITY
- Set appropriate expiry time for JWT, should this be the same as Stravas token expiry?
- Use a proper secret when creating JWTs
- Think about config, It's ok for now but not suitable for production.
- Refresh token if we receive an auth error when communicating with Strava
- Hash access and refresh tokens, don't store as plain text

MONGO
- Research implementing a NoSQL approach - possibly Mongo  ✔️
- Once research is complete on NoSQL approach add a repo to the user service ✔️
- Mongo tests within the user package 
- Add config approach for Mongo, make use of ENV vars for secrets

TESTING
- Add some full service tests using dockertest - look at doing this for handlers package ✔️
- Add a fixture helper when creating dummy objects for test

USERS
- Add some more detail to the user struct, first name and last name (can infer this from the Strava response)
- Finish off Mongo repo tests for users

EVENTS
- Add event struct ✔️
- Add event service ✔️
- Event may be our first aggregate model, think about how to store this
- Add segment entity to fetch a users segments ✔️
- Add Mongo implementation for the event package ✔️
- Query for detailed segments to return a Polyline to display on the FE ✔️
- Add segments to an event ✔️


MISC REFACTORS
- Move the pkg package into a common place for use in Monorepo.
- Add a better read me explaining everything within this repo
- Improve Echo setup within setup package. Look into extracting routes into their own file or package ✔️
- Strava URL & callback URL should come from config
- Make local development possible with the mocked Strava-API. Need an already logged-in user or to mock the
Strava auth ✔️

FE
- Implement basic user flow for creating an event on the FE.
- Design and implement basic registration / login page within FE
- Think about a basic home page, show some stats?
- Error strategy, we currently don't show any server errors on the FE
- Implement React Leaflet to show segments on a map with a Polyline ✔️
- Fix the user authentication 
- Fix null map issues if data isn't present within a component ✔️
- Migrate to use the Fetch API.
