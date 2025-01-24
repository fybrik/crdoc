# API Reference

Packages:

- [project.io/v1alpha1](#projectiov1alpha1)

# project.io/v1alpha1

Resource Types:

- [Example](#example)




## Example
<sup><sup>[↩ Parent](#projectiov1alpha1 )</sup></sup>








<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>project.io/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>Example</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b>arbitrary</b></td>
        <td>map[string]string</td>
        <td>
          arbitrary field<br/>
          <br/>
            <i>Default</i>: map[]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#examplemessagesindexkey">messages</a></b></td>
        <td>[]map[string]object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Example.messages[index][key]
<sup><sup>[↩ Parent](#example)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>area</b></td>
        <td>enum</td>
        <td>
          Some area<br/>
          <br/>
            <i>Enum</i>: East, West<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>code</b></td>
        <td>integer</td>
        <td>
          Some code<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Default</i>: 7<br/>
            <i>Minimum</i>: 0<br/>
            <i>Maximum</i>: 10<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>text</b></td>
        <td>string</td>
        <td>
          Just text<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>zones</b></td>
        <td>[]enum</td>
        <td>
          Some zones<br/>
          <br/>
            <i>Enum</i>: A, B<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
