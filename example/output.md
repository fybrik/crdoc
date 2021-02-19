# API Reference

Packages:

- [project.io/v1alpha1](#project.io/v1alpha1)

# project.io/v1alpha1

Resource Types:

- [Example](#example)




## Example
<sup><sup>[↩ Parent](#project.io/v1alpha1 )</sup></sup>








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
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#examplemessagesindexkey">messages</a></b></td>
        <td>[]map[string]object</td>
        <td></td>
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
        <td>Area... [East West]</td>
        <td>false</td>
      </tr><tr>
        <td><b>code</b></td>
        <td>integer</td>
        <td></td>
        <td>false</td>
      </tr><tr>
        <td><b>text</b></td>
        <td>string</td>
        <td></td>
        <td>false</td>
      </tr></tbody>
</table>