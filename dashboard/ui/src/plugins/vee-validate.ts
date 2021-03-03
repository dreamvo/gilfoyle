import { extend } from "vee-validate";
import {
  email,
  ext,
  integer,
  max,
  max_value as MaxValue,
  min,
  min_value as MinValue,
  regex,
  required,
  size
} from "vee-validate/dist/rules";

extend("required", {
  ...required,
  message: "This field is required"
});

extend("min", {
  ...min,
  message: "Please enter at least {length} characters"
});

extend("max", {
  ...max,
  message: "Please don't enter more than {length} characters"
});

extend("email", {
  ...email,
  message: "Please enter a valid email address"
});

extend("min_value", {
  ...MinValue,
  message: "This value must be more than {min}"
});

extend("max_value", {
  ...MaxValue,
  message: "This value must be less than {max}"
});

extend("confirmedBy", {
  params: ["target"],
  // Target here is the value of the target field
  validate(value: string, { target }: Record<string, any>) {
    return value === target;
  },
  // here it is its name, because we are generating a message
  message: "This field must have the same value as {target}"
});

extend("ext", {
  ...ext,
  message: "That type of file extension is not supported"
});

extend("size", {
  ...size,
  message: "File exceed size limit : {size} kB"
});

extend("integer", {
  ...integer,
  message: "This field must contain an integer"
});

extend("regex", {
  ...regex,
  message: "The value of this field is incorrect"
});
