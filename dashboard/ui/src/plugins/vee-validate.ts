import { extend } from 'vee-validate'
import {
    email,
    min,
    min_value as MinValue,
    max_value as MaxValue,
    required,
    max,
    ext,
    size,
    integer,
    regex,
} from 'vee-validate/dist/rules'

extend('required', {
    ...required,
    message: 'That field is required',
});

extend('min', {
    ...min,
    message: 'Veuillez saisir au moins {length} caractères',
})

extend('max', {
    ...max,
    message: 'Veuillez ne pas dépasser {length} caractères',
})


extend('email', {
    ...email,
    message: 'Veuillez saisir une adresse email valide',
})

extend('min_value', {
    ...MinValue,
    message: 'Ce champs doit être supérieur à {min}',
})

extend('max_value', {
    ...MaxValue,
    message: 'Ce champs doit être inférieur à {max}',
})

extend('confirmedBy', {
    params: ['target'],
    // Target here is the value of the target field
    validate(value, { target }: any) {
        return value === target
    },
    // here it is its name, because we are generating a message
    message: 'Ce champs doit être identique au champs {target}',
})


extend('ext', {
    ...ext,
    message: "Ce type de fichier n'est pas supporté",
})

extend('size', {
    ...size,
    message: 'Ce fichier dépasse la taille maximale autorisée : {size} kB',
})

extend('integer', {
    ...integer,
    message: "Ce champs ne doit contenir qu'un nombre entier",
})

extend('regex', {
    ...regex,
    message: 'Le format de ce champs est incorrect',
})
