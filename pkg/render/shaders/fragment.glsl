#version 410
in vec2 TexCoords;
in vec3 Normal;
in vec3 FragPos;

out vec4 FragColor;

uniform sampler2D moonTexture;

void main() {
    // Light position from camera direction
    vec3 lightPos = vec3(2.0, 2.0, 2.0);
    vec3 lightColor = vec3(1.0);
    
    // Increase ambient light to prevent dark areas
    float ambientStrength = 0.3;
    vec3 ambient = ambientStrength * lightColor;
    
    // Diffuse lighting
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;
    
    // Add minimum light level to prevent complete darkness
    float minLight = 0.3;
    
    // Final color
    vec3 color = texture(moonTexture, TexCoords).rgb;
    vec3 result = (ambient + diffuse + minLight) * color;
    
    // Ensure we don't exceed maximum brightness
    result = min(result, vec3(1.0));
    
    FragColor = vec4(result, 1.0);
}